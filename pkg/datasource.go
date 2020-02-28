package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	ds := &googleSheetsDataSource{
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
		os.Exit(1)
	}
}

type googleSheetsDataSource struct {
	plugin.NetRPCUnsupportedPlugin
	logger      hclog.Logger
	googlesheet *googlesheets.GoogleSheets
}

func (gsd *googleSheetsDataSource) DataQuery(ctx context.Context, req *backend.DataQueryRequest) (*backend.DataQueryResponse, error) {
	res := &backend.DataQueryResponse{}
	config := gs.GoogleSheetConfig{}
	if err := json.Unmarshal(req.PluginConfig.JSONData, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal DataSourceInfo json: %w", err)
	}

	config.ApiKey = req.PluginConfig.DecryptedSecureJSONData["apiKey"]
	config.JWT = req.PluginConfig.DecryptedSecureJSONData["jwt"]

	for _, q := range req.Queries {
		queryModel := &gs.QueryModel{}
		if err := json.Unmarshal(q.JSON, &queryModel); err != nil {
			return nil, fmt.Errorf("failed to unmarshal query: %w", err)
		}

		frame, err := gsd.googlesheet.Query(ctx, q.RefID, queryModel, &config, q.TimeRange)
		if err != nil {
			gsd.logger.Error("Query failed", "refId", q.RefID, "error", err)
			// TEMP: at the moment, the only way to return an error is by using meta
			res.Metadata = map[string]string{"error": err.Error()}
			continue
		}

		res.Frames = append(res.Frames, []*df.Frame{frame}...)
	}

	return res, nil
}

func (gsd *googleSheetsDataSource) CallResource(ctx context.Context, req *backend.CallResourceRequest) (*backend.CallResourceResponse, error) {
	config := gs.GoogleSheetConfig{}
	if err := json.Unmarshal(req.PluginConfig.JSONData, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal configuration: %w", err)
	}

	config.ApiKey = req.PluginConfig.DecryptedSecureJSONData["apiKey"]
	config.JWT, _ = req.PluginConfig.DecryptedSecureJSONData["jwt"]

	response := make(map[string]interface{})
	var res interface{}
	var err error
	switch req.Path {
	case "spreadsheets":
		res, err = gsd.googlesheet.GetSpreadsheets(ctx, &config)
	case "test":
		err = gsd.googlesheet.TestAPI(ctx, &config)
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
