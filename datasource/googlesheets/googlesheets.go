package googlesheets

import (
	"fmt"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	cd "github.com/grafana/google-sheets-datasource/datasource/columndefinition"
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
	Cache *cache.Cache
	Logger      hclog.Logger
}

// Query function
func (gs *GoogleSheets) Query(ctx context.Context, refID string, qm *QueryModel, config *GoogleSheetConfig, timeRange backend.TimeRange) ([]*df.Frame, error) {
	srv, err := createSheetsService(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Unable to create service: %v", err.Error())
	}

	sheet, err := gs.getSpreadSheet(srv, qm, config)
	if err != nil {
		return nil, err
	}

	return gs.transformSheetToDataFrames(sheet, refID, qm)
}


// TestAPI function
func (gs *GoogleSheets) TestAPI(ctx context.Context, config *GoogleSheetConfig) ([]*df.Frame, error) {
	_, err := createSheetsService(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Invalid datasource configuration: %s", err)
	}

	return []*df.Frame{df.New("TestAPI")}, nil
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

func (gs *GoogleSheets) getSpreadSheet(srv *sheets.Service, qm *QueryModel, config *GoogleSheetConfig) (*sheets.GridData, error) {
	var sheet *sheets.GridData
	cacheKey:=qm.SpreadsheetID + qm.Range
	if item, found := gs.Cache.Get(cacheKey); found {
		sheet = item.(*sheets.GridData)
		gs.Logger.Debug("GET_FROM_CACHE: ", cacheKey)
	} else {
		result, err := srv.Spreadsheets.Get(qm.SpreadsheetID).Ranges(qm.Range).IncludeGridData(true).Do()
		if err != nil {
			return nil, fmt.Errorf("Unable to get spreadsheet: %v", err.Error())
		} else if config.CacheDurationSeconds > 0 {
			gs.Logger.Debug("PUT_IN_CACHE: ", cacheKey)
			sheet = result.Sheets[0].Data[0]
			gs.Cache.Set(cacheKey, sheet, time.Duration(config.CacheDurationSeconds)*time.Second)
		}
	}

	return sheet, nil
}

func (gs *GoogleSheets) transformSheetToDataFrames(sheet *sheets.GridData, refID string, qm *QueryModel) ([]*df.Frame, error) {
	fields := []*df.Field{}
	columns := getColumnDefintions(sheet.RowData)

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

		field.Config = &df.FieldConfig{}
		field.Config.Unit = column.GetUnit()

		if column.HasMixedTypes() {
			gs.Logger.Error("Multipe data types found in column " + column.Header + ". Using string data type")
		}

		if column.HasMixedUnits() {
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
					return []*df.Frame{frame}, fmt.Errorf("error while parsing date : %v", err.Error())
				}
				frame.Fields[columnIndex].Vector.Set(rowIndex-1, &time)
			case "NUMBER":
				if cellData.EffectiveValue != nil {
					frame.Fields[columnIndex].Vector.Set(rowIndex-1, &cellData.EffectiveValue.NumberValue)
				}
			case "STRING":
				frame.Fields[columnIndex].Vector.Set(rowIndex-1, &cellData.FormattedValue)
			}
		}
	}

	return []*df.Frame{frame}, nil
}

func createSheetsService(ctx context.Context, config *GoogleSheetConfig) (*sheets.Service, error) {
	if config.AuthType == "none" {
		return sheets.NewService(ctx, option.WithAPIKey(config.ApiKey))
	}

	jwtConfig, err := google.JWTConfigFromJSON([]byte(config.JwtFile), "https://www.googleapis.com/auth/spreadsheets.readonly", "https://www.googleapis.com/auth/spreadsheets")
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

	jwtConfig, err := google.JWTConfigFromJSON([]byte(config.JwtFile), drive.DriveMetadataReadonlyScope)
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