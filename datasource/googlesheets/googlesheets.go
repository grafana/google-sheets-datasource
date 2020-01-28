package googlesheets

import (
	df "github.com/grafana/grafana-plugin-sdk-go/dataframe"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Query function
func Query(ctx context.Context, refID string, sheet GoogleSheetRangeInfo, config GoogleSheetConfig) (*df.Frame, error) {
	srv, err := sheets.NewService(ctx, option.WithAPIKey(config.ApiKey))
	if err != nil {
		// gsd.logger.Error("Unable to create service: %v", err.Error())
		return nil, err
	}

	resp, err := srv.Spreadsheets.Values.Get(sheet.SpreadsheetID, sheet.Range).Do()
	if err != nil {
		// gsd.logger.Error("Unable to retrieve data from sheet: %v", err.Error())
		return nil, err
	}

	fields := []*df.Field{}
	for _, column := range resp.Values[0] {
		fields = append(fields, df.NewField(column.(string), nil, []string{}))
	}

	frame := df.New(sheet.Range, fields...)
	frame.RefID = refID

	for index := 1; index < len(resp.Values); index++ {
		for columnID, value := range resp.Values[index] {
			frame.Fields[columnID].Vector.Append(value.(string))
		}
	}

	return frame, nil
}
