package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/google-sheets-datasource/pkg/googlesheets"
	gs "github.com/grafana/google-sheets-datasource/pkg/googlesheets"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	df "github.com/grafana/grafana-plugin-sdk-go/dataframe"
	"github.com/patrickmn/go-cache"

	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"

	"golang.org/x/net/context"
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
	cache := cache.New(300*time.Second, 5*time.Second)
	ds := &GoogleSheetsDataSource{
		logger: pluginLogger,
		googlesheet: &googlesheets.GoogleSheets{
			Cache:  cache,
			Logger: pluginLogger,
		},
	}
	err := backend.Serve(backend.ServeOpts{
		DataQueryHandler: ds,
	})
	if err != nil {
		pluginLogger.Error(err.Error())
	}
}

type GoogleSheetsDataSource struct {
	plugin.NetRPCUnsupportedPlugin
	logger      hclog.Logger
	googlesheet *googlesheets.GoogleSheets
}

func (gsd *GoogleSheetsDataSource) DataQuery(ctx context.Context, req *backend.DataQueryRequest) (*backend.DataQueryResponse, error) {
	res := &backend.DataQueryResponse{}
	config := gs.GoogleSheetConfig{}
	err := json.Unmarshal(req.PluginConfig.JSONData, &config)
	if err != nil {
		gsd.logger.Error("Could not unmarshal DataSourceInfo json", "Error", err)
		return nil, err
	}

	for _, q := range req.Queries {
		queryModel := &gs.QueryModel{}
		err := json.Unmarshal(q.JSON, &queryModel)

		if err != nil {
			gsd.logger.Error("Failed to unmarshal query: %v", err.Error())
			return nil, fmt.Errorf("Invalid query")
		}

		var frames []*df.Frame
		switch queryModel.QueryType {
		case "testAPI":
			frames, err = gsd.googlesheet.TestAPI(ctx, &config)
		case "query":
			frames, err = gsd.googlesheet.Query(ctx, q.RefID, queryModel, &config, q.TimeRange)
		case "getSpreadsheets":
			spreadSheets, _ := gsd.googlesheet.GetSpreadsheetsByServiceAccount(ctx, &config)
			frame := df.New("getHeader")
			res.Metadata = spreadSheets
			frames = []*df.Frame{frame}
		default:
			return nil, fmt.Errorf("Invalid query type")
		}

		if err != nil {
			gsd.logger.Debug("Metric Error: ", err.Error())
		}

		// if err != nil {
		// 	gsd.logger.Debug("QueryError", "QueryError", err.Error())
		// 	frame = df.New("default")
		// 	frame.RefID = q.RefID
		// 	frame.Meta = &df.QueryResultMeta{Custom: make(map[string]interface{})}
		// 	frame.Meta.Custom["error"] = err.Error()
		// 	// return nil, err
		// }
		res.Frames = append(res.Frames, frames...)
	}

	return res, nil
}
