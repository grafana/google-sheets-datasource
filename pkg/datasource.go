package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/grafana/google-sheets-datasource/pkg/googlesheets"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	df "github.com/grafana/grafana-plugin-sdk-go/dataframe"
	"github.com/patrickmn/go-cache"

	hclog "github.com/hashicorp/go-hclog"

	"context"
)

const (
	pluginID = "google-sheets-datasource"
)

func main() {
	pluginLogger := hclog.New(&hclog.LoggerOptions{
		Name: pluginID,
		// TODO: How to make level configurable?
		Level: hclog.LevelFromString("DEBUG"),
	})
	cache := cache.New(300*time.Second, 5*time.Second)
	ds := &googleSheetsDataSource{
		logger: pluginLogger,
		googlesheet: &googlesheets.GoogleSheets{
			Cache:  cache,
			Logger: pluginLogger,
		},
	}
	if err := backend.Serve(backend.ServeOpts{
		DataQueryHandler:    ds,
		CallResourceHandler: ds,
	}); err != nil {
		pluginLogger.Error(err.Error())
		os.Exit(1)
	}
}

type googleSheetsDataSource struct {
	logger      hclog.Logger
	googlesheet *googlesheets.GoogleSheets
}

// DataQuery queries for data.
func (gsd *googleSheetsDataSource) DataQuery(ctx context.Context, req *backend.DataQueryRequest) (*backend.DataQueryResponse, error) {
	res := &backend.DataQueryResponse{}
	config := googlesheets.GoogleSheetConfig{}
	if err := json.Unmarshal(req.PluginConfig.JSONData, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal DataSourceInfo json: %w", err)
	}

	config.APIKey = req.PluginConfig.DecryptedSecureJSONData["apiKey"]
	config.JWT = req.PluginConfig.DecryptedSecureJSONData["jwt"]

	for _, q := range req.Queries {
		queryModel := &googlesheets.QueryModel{}
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

// CallResource calls a resource.
func (gsd *googleSheetsDataSource) CallResource(ctx context.Context, req *backend.CallResourceRequest) (*backend.CallResourceResponse, error) {
	config := googlesheets.GoogleSheetConfig{}
	if err := json.Unmarshal(req.PluginConfig.JSONData, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal configuration: %w", err)
	}

	config.APIKey = req.PluginConfig.DecryptedSecureJSONData["apiKey"]
	config.JWT = req.PluginConfig.DecryptedSecureJSONData["jwt"]

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
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return &backend.CallResourceResponse{
		Status: http.StatusOK,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: body,
	}, nil
}
