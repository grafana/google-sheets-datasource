package googlesheets

import (
	"fmt"
	"time"

	// "github.com/davecgh/go-spew/spew"
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

func getTypeDefaultValue(t string) interface{} {
	switch t {
	case "time":
		return nil
	case "float64":
		return 0.0
	default:
		return ""
	}
}

func getTableData(srv *sheets.Service, refID string, qm *QueryModel, logger hclog.Logger) ([]*df.Frame, error) {
	resp, err := srv.Spreadsheets.Values.Get(qm.SpreadsheetID, qm.Range).MajorDimension(qm.MajorDimension).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve data from sheet: %v", err.Error())
	}

	fields := []*df.Field{}
	for _, column := range resp.Values[0] {
		fields = append(fields, df.NewField(column.(string), nil, []string{}))
	}

	frame := df.New(qm.Range, fields...)
	frame.RefID = refID

	for index := 1; index < len(resp.Values); index++ {
		for columnID := 0; columnID < len(fields); columnID++ {
			if columnID+1 <= len(resp.Values[index]) {
				frame.Fields[columnID].Vector.Append(resp.Values[index][columnID].(string))
			} else {
				logger.Debug("appending default value", string(columnID))
				frame.Fields[columnID].Vector.Append("")
			}

		}
	}

	return []*df.Frame{frame}, nil
}

func getTimeSeriesData(srv *sheets.Service, refID string, qm *QueryModel, timeRange backend.TimeRange, logger hclog.Logger) ([]*df.Frame, error) {
	result, err := srv.Spreadsheets.Get(qm.SpreadsheetID).IncludeGridData(true).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to get spreadsheet: %v", err.Error())
	}

	sheet := result.Sheets[0].Data[0]
	if result.Properties.TimeZone != "" {
		loc, err := time.LoadLocation(result.Properties.TimeZone)
		if err != nil {
			return nil, fmt.Errorf("error while loading timezone: ", err.Error())
		}
		time.Local = loc
	}

	frames := []*df.Frame{}
	for _, metric := range qm.MetricColumns {
		frame := df.New(metric.Label,
			df.NewField("Time", nil, []time.Time{}),
			df.NewField("Value", nil, []float64{}),
		)

		frame.RefID = refID
		frames = append(frames, frame)
	}

	for rowIndex := 1; rowIndex < len(sheet.RowData); rowIndex++ {
		timeCell := sheet.RowData[rowIndex].Values[qm.TimeColumn.Value]
		time, err := dateparse.ParseLocal(timeCell.FormattedValue)
		if err != nil {
			logger.Error("error while parsing date :", err.Error())
			continue
		}

		if time.Sub(timeRange.From) < 0 || time.Sub(timeRange.To) > 0 {
			logger.Debug("time out of time range")
			continue
		}

		for i, metric := range qm.MetricColumns {
			frames[i].Fields[0].Vector.Append(time)

			if metric.Value+1 > len(sheet.RowData[rowIndex].Values) {
				frames[i].Fields[1].Vector.Append(0.0)
			} else {
				metricCell := sheet.RowData[rowIndex].Values[metric.Value]
				if metricCell.EffectiveValue == nil {
					frames[i].Fields[1].Vector.Append(0.0)
				} else {
					frames[i].Fields[1].Vector.Append(metricCell.EffectiveValue.NumberValue)
				}
			}
		}
	}

	return frames, nil
}

// Query function
func Query(ctx context.Context, refID string, qm *QueryModel, config *GoogleSheetConfig, timeRange backend.TimeRange, logger hclog.Logger) ([]*df.Frame, error) {
	srv, err := createSheetsService(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Unable to create service: %v", err.Error())
	}

	switch qm.ResultFormat {
	case "TABLE":
		return getTableData(srv, refID, qm, logger)
	case "TIME_SERIES":
		return getTimeSeriesData(srv, refID, qm, timeRange, logger)
	default:
		return nil, fmt.Errorf("Invalid result format: %v", qm.ResultFormat)
	}

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
		logger.Debug("aaaaa12: " + err.Error())
		return nil, fmt.Errorf("Could not get all files: %s", err.Error())
	}
	fileNames := map[string]string{}
	for _, i := range files {
		logger.Debug("%s (%s)\n", i.Name, i.Id)
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
