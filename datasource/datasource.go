package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	df "github.com/grafana/grafana-plugin-sdk-go/dataframe"

	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"

	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	pluginID     = "google-sheets-datasource"
	variableName = "GOOGLE_SHEETS_API_KEY"
)

var pluginLogger = hclog.New(&hclog.LoggerOptions{
	Name:  pluginID,
	Level: hclog.LevelFromString("DEBUG"),
})

func main() {
	path, exist := os.LookupEnv(variableName)
	if !exist {
		pluginLogger.Error("could not read environment variable", variableName)
		panic(fmt.Errorf("could not read environment variable %v", variableName))
	} else {
		pluginLogger.Debug("environment variable for google sheets found", "variable", variableName, "value", path)
	}

	err := backend.Serve(backend.ServeOpts{
		DataQueryHandler: &GoogleSheetsDataSource{
			logger: pluginLogger,
		},
	})
	if err != nil {
		pluginLogger.Error(err.Error())
	}
}

type GoogleSheetsDataSource struct {
	plugin.NetRPCUnsupportedPlugin
	logger hclog.Logger
}

type GoogleSheetRangeInfo struct {
	QueryType     string
	SpreadsheetID string
	Range         string
}

func (gsd *GoogleSheetsDataSource) DataQuery(ctx context.Context, req *backend.DataQueryRequest) (*backend.DataQueryResponse, error) {
	res := &backend.DataQueryResponse{}

	for _, q := range req.Queries {
		googleSheetRangeInfo := &GoogleSheetRangeInfo{}
		err := json.Unmarshal(q.JSON, &googleSheetRangeInfo)

		apiKey, _ := os.LookupEnv(variableName)
		srv, err := sheets.NewService(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			gsd.logger.Error("Unable to create service: %v", err.Error())
		}
		resp, err := srv.Spreadsheets.Values.Get(googleSheetRangeInfo.SpreadsheetID, googleSheetRangeInfo.Range).Do()
		if err != nil {
			gsd.logger.Error("Unable to retrieve data from sheet: %v", err.Error())
		}

		fields := []*df.Field{}
		for _, column := range resp.Values[0] {
			fields = append(fields, df.NewField(column.(string), nil, []string{}))
		}

		frame := df.New(googleSheetRangeInfo.Range, fields...)
		frame.RefID = q.RefID

		for index := 1; index < len(resp.Values); index++ {
			for columnID, value := range resp.Values[index] {
				frame.Fields[columnID].Vector.Append(value.(string))
			}
		}

		res.Frames = append(res.Frames, frame)
	}

	return res, nil
}
