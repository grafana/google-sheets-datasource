package googlesheets

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	df "github.com/grafana/grafana-plugin-sdk-go/dataframe"
	"github.com/hashicorp/go-hclog"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

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

func getTableData(srv *sheets.Service, refID string, qm *QueryModel, logger hclog.Logger) ([]*df.Frame, error) {
	result, err := srv.Spreadsheets.Get(qm.SpreadsheetID).Ranges(qm.Range).IncludeGridData(true).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to get spreadsheet: %v", err.Error())
	}

	sheet := result.Sheets[0].Data[0]

	fields := []*df.Field{}
	columns := []string{}
	for columnIndex, column := range sheet.RowData[0].Values {
		// columnTypes := getColumnTypes(sheet.RowData)
		// for columnIndex, columnType := range columnTypes {
		// 	logger.Debug("COLUMN TYPE", spew.Sdump(columnType, columnIndex))
		// }

		columnType := getColumnType(columnIndex, sheet.RowData)
		columns = append(columns, columnType)
		switch columnType {
		case "TIME":
			fields = append(fields, df.NewField(column.FormattedValue, nil, []*time.Time{}))
		case "NUMBER":
			fields = append(fields, df.NewField(column.FormattedValue, nil, []*float64{}))
		case "STRING":
			fields = append(fields, df.NewField(column.FormattedValue, nil, []*string{}))
		}

		fields[columnIndex].Config = &df.FieldConfig{Unit: getColumnUnit(columnIndex, sheet.RowData)}
	}

	frame := df.New(refID,
		fields...,
	)

	for rowIndex := 1; rowIndex < len(sheet.RowData); rowIndex++ {
		for columnIndex, columnType := range columns {
			if columnIndex < len(sheet.RowData[rowIndex].Values) {
				cellData := sheet.RowData[rowIndex].Values[columnIndex]
				switch columnType {
				case "TIME":
					time, err := dateparse.ParseLocal(cellData.FormattedValue)
					if err != nil {
						return []*df.Frame{frame}, fmt.Errorf("error while parsing date :", err.Error())
					}
					frame.Fields[columnIndex].Vector.Append(&time)
				case "NUMBER":
					frame.Fields[columnIndex].Vector.Append(&cellData.EffectiveValue.NumberValue)
				case "STRING":
					frame.Fields[columnIndex].Vector.Append(&cellData.FormattedValue)
				}
			}
		}
	}

	return []*df.Frame{frame}, nil
}

// Query function
func Query(ctx context.Context, refID string, qm *QueryModel, config *GoogleSheetConfig, timeRange backend.TimeRange, logger hclog.Logger) ([]*df.Frame, error) {
	srv, err := createSheetsService(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Unable to create service: %v", err.Error())
	}

	return getTableData(srv, refID, qm, logger)
}

func allFiles(d *drive.Service) ([]*drive.File, error) {
	var fs []*drive.File
	pageToken := ""
	for {
		q := d.Files.List().Q("mimeType='application/vnd.google-apps.spreadsheet'")
		// If we have a pageToken set, apply it to the query
		if pageToken != "" {
			q = q.PageToken(pageToken)
		}
		r, err := q.Do()
		if err != nil {
			return fs, err
		}
		//   fs = append(fs, r.Items...)
		fs = append(fs, r.Files...)
		pageToken = r.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return fs, nil
}

//GetSpreadsheets
func GetSpreadsheets(ctx context.Context, config *GoogleSheetConfig, logger hclog.Logger) (map[string]string, error) {
	srv, err := createDriveService(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Invalid datasource configuration: %s", err)
	}

	files, err := allFiles(srv)
	if err != nil {
		return nil, fmt.Errorf("Could not get all files: %s", err.Error())
	}
	fileNames := map[string]string{}
	for _, i := range files {
		fileNames[i.Id] = i.Name
	}

	return fileNames, nil
}

func GetHeaders(ctx context.Context, qm *QueryModel, config *GoogleSheetConfig, logger hclog.Logger) ([]string, error) {
	srv, err := createSheetsService(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Invalid datasource configuration: %s", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(qm.SpreadsheetID, qm.Range).MajorDimension(qm.MajorDimension).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve data from sheet: %v", err.Error())
	}

	headers := []string{}
	for _, column := range resp.Values[0] {
		headers = append(headers, column.(string))
	}
	return headers, nil
}

// TestAPI function
func TestAPI(ctx context.Context, config *GoogleSheetConfig) ([]*df.Frame, error) {
	_, err := createSheetsService(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Invalid datasource configuration: %s", err)
	}

	return []*df.Frame{df.New("TestAPI")}, nil
}
