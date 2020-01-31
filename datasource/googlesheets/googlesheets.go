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
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func createService(ctx context.Context, config *GoogleSheetConfig) (*sheets.Service, error) {
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

func getTableData(srv *sheets.Service, refID string, qm *QueryModel, logger hclog.Logger) (*df.Frame, error) {
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
		for columnID, value := range resp.Values[index] {
			frame.Fields[columnID].Vector.Append(value.(string))
		}
	}

	return frame, nil
}

func getTimeSeriesData(srv *sheets.Service, refID string, qm *QueryModel, timeRange backend.TimeRange, logger hclog.Logger) (*df.Frame, error) {
	result, err := srv.Spreadsheets.Get(qm.SpreadsheetID).Ranges(qm.Range).IncludeGridData(true).Do()
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

	frame := df.New(qm.Range,
		df.NewField(qm.TimeColumn.Label, nil, []time.Time{}),
		df.NewField(qm.MetricColumn.Label, nil, []float64{}),
	)

	for rowIndex := 1; rowIndex < len(sheet.RowData); rowIndex++ {
		timeCell := sheet.RowData[rowIndex].Values[qm.TimeColumn.Value]
		time, err := dateparse.ParseLocal(timeCell.FormattedValue)
		if err != nil {
			return nil, fmt.Errorf("error while parsing date :", err.Error())
		}

		if time.Sub(timeRange.From) < 0 {
			logger.Debug("before time range")
			continue
		}

		if time.Sub(timeRange.To) > 0 {
			logger.Debug("after time range")
			continue
		}
		frame.Fields[qm.TimeColumn.Value].Vector.Append(time)

		metricCell := sheet.RowData[rowIndex].Values[qm.MetricColumn.Value]
		frame.Fields[qm.MetricColumn.Value].Vector.Append(metricCell.EffectiveValue.NumberValue)
	}

	return frame, nil
}

// Query function
func Query(ctx context.Context, refID string, qm *QueryModel, config *GoogleSheetConfig, timeRange backend.TimeRange, logger hclog.Logger) (*df.Frame, error) {
	srv, err := createService(ctx, config)
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

func GetHeaders(ctx context.Context, qm *QueryModel, config *GoogleSheetConfig, logger hclog.Logger) ([]string, error) {
	srv, err := createService(ctx, config)
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
func TestAPI(ctx context.Context, config *GoogleSheetConfig) (*df.Frame, error) {
	_, err := createService(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Invalid datasource configuration: %s", err)
	}

	return df.New("TestAPI"), nil
}
