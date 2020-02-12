package googlesheets

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/davecgh/go-spew/spew"
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
	columns := getColumnDefintions(sheet.RowData, logger)

	for _, columnDef := range columns {
		logger.Debug("COLUMNNAME", spew.Sdump(columnDef.Header))
		var field *df.Field
		switch columnDef.Type {
		case "TIME":
			field = df.NewField(columnDef.Header, nil, []*time.Time{})
		case "NUMBER":
			field = df.NewField(columnDef.Header, nil, []*float64{})
		case "STRING":
			field = df.NewField(columnDef.Header, nil, []*string{})
		}

		if columnDef.Unit != "" {
			field.Config = &df.FieldConfig{Unit: columnDef.Unit}
		}

		fields = append(fields, field)
	}
	logger.Debug("COLUMN fields", spew.Sdump(fields))

	frame := df.New(refID,
		fields...,
	)

	for rowIndex := 1; rowIndex < len(sheet.RowData); rowIndex++ {
		for _, columnDef := range columns {
			if columnDef.ColumnIndex < len(sheet.RowData[rowIndex].Values) {
				cellData := sheet.RowData[rowIndex].Values[columnDef.ColumnIndex]
				switch columnDef.Type {
				case "TIME":
					time, err := dateparse.ParseLocal(cellData.FormattedValue)
					if err != nil {
						return []*df.Frame{frame}, fmt.Errorf("error while parsing date :", err.Error())
					}
					frame.Fields[columnDef.ColumnIndex].Vector.Append(&time)
				case "NUMBER":
					frame.Fields[columnDef.ColumnIndex].Vector.Append(&cellData.EffectiveValue.NumberValue)
				case "STRING":
					frame.Fields[columnDef.ColumnIndex].Vector.Append(&cellData.FormattedValue)
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

// TestAPI function
func TestAPI(ctx context.Context, config *GoogleSheetConfig) ([]*df.Frame, error) {
	_, err := createSheetsService(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Invalid datasource configuration: %s", err)
	}

	return []*df.Frame{df.New("TestAPI")}, nil
}
