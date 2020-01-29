package googlesheets

import (
	"fmt"

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

// Query function
func Query(ctx context.Context, refID string, sheet *QueryModel, config *GoogleSheetConfig, logger hclog.Logger) (*df.Frame, error) {
	srv, err := createService(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Unable to create service: %v", err.Error())
	}

	resp, err := srv.Spreadsheets.Values.Get(sheet.SpreadsheetID, sheet.Range).MajorDimension(sheet.MajorDimension).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve data from sheet: %v", err.Error())
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

// TestAPI function
func TestAPI(ctx context.Context, config *GoogleSheetConfig) (*df.Frame, error) {
	_, err := createService(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Invalid datasource configuration: %s", err)
	}

	return df.New("TestAPI"), nil
}
