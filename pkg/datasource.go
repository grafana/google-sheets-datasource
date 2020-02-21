package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		DataQueryHandler:    ds,
		CallResourceHandler: ds,
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

	config.ApiKey = req.PluginConfig.DecryptedSecureJSONData["apiKey"]
	config.JWT, _ = req.PluginConfig.DecryptedSecureJSONData["jwt"]

	for _, q := range req.Queries {
		queryModel := &gs.QueryModel{}
		err := json.Unmarshal(q.JSON, &queryModel)

		if err != nil {
			gsd.logger.Error("Failed to unmarshal query: %v", err.Error())
			return nil, fmt.Errorf("Invalid query")
		}

		var frame *df.Frame
		switch queryModel.QueryType {
		case "testAPI":
			frame, err = gsd.googlesheet.TestAPI(ctx, &config)
		case "query":
			frame, err = gsd.googlesheet.Query(ctx, q.RefID, queryModel, &config, q.TimeRange)
		default:
			return nil, fmt.Errorf("Invalid query type")
		}

		if err != nil {
			// TEMP: at the moment, the only way to return an error is by using meta
			res.Metadata = map[string]string{"error": err.Error()}
		}

		res.Frames = append(res.Frames, []*df.Frame{frame}...)
	}

	return res, nil
}

func (gsd *GoogleSheetsDataSource) CallResource(ctx context.Context, req *backend.CallResourceRequest) (*backend.CallResourceResponse, error) {
	config := gs.GoogleSheetConfig{}
	err := json.Unmarshal(req.PluginConfig.JSONData, &config)
	if err != nil {
		gsd.logger.Error("Could not unmarshal DataSourceInfo json", "Error", err)
		return nil, err
	}
	config.ApiKey = req.PluginConfig.DecryptedSecureJSONData["apiKey"]
	config.JWT, _ = req.PluginConfig.DecryptedSecureJSONData["jwt"]

	response := make(map[string]interface{})
	var res interface{}
	switch req.Path {
	case "spreadsheets":
		res, err = gsd.googlesheet.GetSpreadsheetsByServiceAccount(ctx, &config)
	}

	if err != nil {
		response["error"] = err.Error()
	} else {
		response[req.Path] = res
	}

	body, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	headers := make(http.Header)
	headers.Add("Content-Type", "application/json")

	return &backend.CallResourceResponse{
		Status:  http.StatusOK,
		Headers: headers,
		Body:    body,
	}, nil
}
