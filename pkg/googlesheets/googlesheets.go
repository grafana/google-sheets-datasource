package googlesheets

import (
	"fmt"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/davecgh/go-spew/spew"
	cd "github.com/grafana/google-sheets-datasource/pkg/googlesheets/columndefinition"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	df "github.com/grafana/grafana-plugin-sdk-go/dataframe"
	"github.com/hashicorp/go-hclog"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GoogleSheets struct {
	Cache  *cache.Cache
	Logger hclog.Logger
}

// Query function
func (gs *GoogleSheets) Query(ctx context.Context, refID string, qm *QueryModel, config *GoogleSheetConfig, timeRange backend.TimeRange) (*df.Frame, error) {
	srv, err := createSheetsService(ctx, config)
	if err != nil {
		return df.New(refID), fmt.Errorf("Unable to create service: %v", err.Error())
	}

	sheet, meta, err := gs.getSpreadSheet(srv, qm, config)
	if err != nil {
		return df.New(refID), err
	}

	frame, err := gs.transformSheetToDataFrame(sheet, meta, refID, qm)
	if err != nil {
		return df.New(refID), err
	}

	return frame, nil
}

// TestAPI function
func (gs *GoogleSheets) TestAPI(ctx context.Context, config *GoogleSheetConfig) (*df.Frame, error) {
	_, err := createSheetsService(ctx, config)
	return df.New("TestAPI"), err
}

//GetSpreadsheets
func (gs *GoogleSheets) GetSpreadsheetsByServiceAccount(ctx context.Context, config *GoogleSheetConfig) (map[string]string, error) {
	srv, err := createDriveService(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Invalid datasource configuration: %s", err)
	}

	files, err := getAllSheets(srv)
	if err != nil {
		return nil, fmt.Errorf("Could not get all files: %s", err.Error())
	}
	fileNames := map[string]string{}
	for _, i := range files {
		fileNames[i.Id] = i.Name
	}

	return fileNames, nil
}

func (gs *GoogleSheets) getSpreadSheet(srv *sheets.Service, qm *QueryModel, config *GoogleSheetConfig) (*sheets.GridData, map[string]interface{}, error) {
	cacheKey := qm.Spreadsheet.ID + qm.Range
	if item, expires, found := gs.Cache.GetWithExpiration(cacheKey); found && qm.CacheDurationSeconds > 0 {
		return item.(*sheets.GridData), map[string]interface{}{"hit": true, "count": gs.Cache.ItemCount(), "expires": fmt.Sprintf("%vs", int(expires.Sub(time.Now()).Seconds()))}, nil
	}

	result, err := srv.Spreadsheets.Get(qm.Spreadsheet.ID).Ranges(qm.Range).IncludeGridData(true).Do()
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to get spreadsheet: %v", err.Error())
	}
	sheet := result.Sheets[0].Data[0]
	gs.Cache.Set(cacheKey, sheet, time.Duration(qm.CacheDurationSeconds)*time.Second)

	return sheet, map[string]interface{}{"hit": false}, nil
}

func (gs *GoogleSheets) transformSheetToDataFrame(sheet *sheets.GridData, meta map[string]interface{}, refID string, qm *QueryModel) (*df.Frame, error) {
	fields := []*df.Field{}
	columns := getColumnDefintions(sheet.RowData)
	warnings := []string{}

	for _, column := range columns {

		var field *df.Field
		switch column.GetType() {
		case "TIME":
			field = df.NewField(column.Header, nil, make([]*time.Time, len(sheet.RowData)-1))
		case "NUMBER":
			field = df.NewField(column.Header, nil, make([]*float64, len(sheet.RowData)-1))
		case "STRING":
			field = df.NewField(column.Header, nil, make([]*string, len(sheet.RowData)-1))
		}

		field.Config = &df.FieldConfig{Unit: column.GetUnit()}

		if column.HasMixedTypes() {
			warnings = append(warnings, fmt.Sprintf("Multipe data types found in column %s. Using string data type", column.Header))
			gs.Logger.Error("Multipe data types found in column " + column.Header + ". Using string")
		}

		if column.HasMixedUnits() {
			warnings = append(warnings, fmt.Sprintf("Multipe units found in column %s. Formatted value will be used", column.Header))
			gs.Logger.Error("Multipe units found in column " + column.Header + ". Formatted value will be used")
		}

		fields = append(fields, field)
	}

	frame := df.New(refID,
		fields...,
	)

	for rowIndex := 1; rowIndex < len(sheet.RowData); rowIndex++ {
		for columnIndex, cellData := range sheet.RowData[rowIndex].Values {
			switch columns[columnIndex].GetType() {
			case "TIME":
				time, err := dateparse.ParseLocal(cellData.FormattedValue)
				if err != nil {
					warnings = append(warnings, fmt.Sprintf("Error while parsing date at row %v in column %s", rowIndex+1, columns[columnIndex].Header))
				} else {
					frame.Fields[columnIndex].Vector.Set(rowIndex-1, &time)
				}
			case "NUMBER":
				if cellData.EffectiveValue != nil {
					frame.Fields[columnIndex].Vector.Set(rowIndex-1, &cellData.EffectiveValue.NumberValue)
				}
			case "STRING":
				frame.Fields[columnIndex].Vector.Set(rowIndex-1, &cellData.FormattedValue)
			}
		}
	}

	meta["warnings"] = warnings
	meta["spreadsheetId"] = qm.Spreadsheet.ID
	meta["range"] = qm.Range
	frame.Meta = &df.QueryResultMeta{Custom: meta}
	gs.Logger.Debug("frame.Meta", spew.Sdump(frame.Meta))

	return frame, nil
}

func createSheetsService(ctx context.Context, config *GoogleSheetConfig) (*sheets.Service, error) {
	if config.AuthType == "none" {
		if len(config.ApiKey) == 0 {
			return nil, fmt.Errorf("Invalid API Key")
		}

		return sheets.NewService(ctx, option.WithAPIKey(config.ApiKey))
	}

	jwtConfig, err := google.JWTConfigFromJSON([]byte(config.JWT), "https://www.googleapis.com/auth/spreadsheets.readonly", "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, fmt.Errorf("Error parsin JWT file: %v", err)
	}

	client := jwtConfig.Client(ctx)
	return sheets.New(client)
}

func createDriveService(ctx context.Context, config *GoogleSheetConfig) (*drive.Service, error) {
	if config.AuthType == "none" {
		return drive.NewService(ctx, option.WithAPIKey(config.ApiKey))
	}

	jwtConfig, err := google.JWTConfigFromJSON([]byte(config.JWT), drive.DriveMetadataReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("Error parsin JWT file: %v", err)
	}

	client := jwtConfig.Client(ctx)

	return drive.New(client)
}

func getColumnDefintions(rows []*sheets.RowData) []*cd.ColumnDefinition {
	columnTypes := []*cd.ColumnDefinition{}
	headerRow := rows[0].Values

	for columnIndex, headerCell := range headerRow {
		columnTypes = append(columnTypes, cd.New(strings.TrimSpace(headerCell.FormattedValue), columnIndex))
	}

	for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
		for columnIndex, columnCell := range rows[rowIndex].Values {
			columnTypes[columnIndex].CheckCell(columnCell)
		}
	}

	return columnTypes
}

func getAllSheets(d *drive.Service) ([]*drive.File, error) {
	var fs []*drive.File
	pageToken := ""
	for {
		q := d.Files.List().Q("mimeType='application/vnd.google-apps.spreadsheet'")
		if pageToken != "" {
			q = q.PageToken(pageToken)
		}
		r, err := q.Do()
		if err != nil {
			return fs, err
		}
		fs = append(fs, r.Files...)
		pageToken = r.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return fs, nil
}
