package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/grafana/google-sheets-datasource/pkg/googlesheets"
	"github.com/grafana/google-sheets-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"

	"context"
)

const metricNamespace = "sheets_datasource"

// GoogleSheetsDataSource handler for google sheets
type GoogleSheetsDataSource struct {
	googlesheet *googlesheets.GoogleSheets
}

// NewDataSource creates the google sheets datasource and sets up all the routes
func NewDataSource(mux *http.ServeMux) *GoogleSheetsDataSource {
	queriesTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "data_query_total",
			Help:      "data query counter",
			Namespace: metricNamespace,
		},
		[]string{"scenario"},
	)
	prometheus.MustRegister(queriesTotal)

	cache := cache.New(300*time.Second, 5*time.Second)
	ds := &GoogleSheetsDataSource{
		googlesheet: &googlesheets.GoogleSheets{
			Cache: cache,
		},
	}

	mux.HandleFunc("/spreadsheets", ds.handleResourceSpreadsheets)
	return ds
}

func readConfig(pluginConfig backend.PluginConfig) (*models.GoogleSheetConfig, error) {
	config := models.GoogleSheetConfig{}
	if err := json.Unmarshal(pluginConfig.DataSourceConfig.JSONData, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal DataSourceInfo json: %w", err)
	}
	config.APIKey = pluginConfig.DataSourceConfig.DecryptedSecureJSONData["apiKey"]
	config.JWT = pluginConfig.DataSourceConfig.DecryptedSecureJSONData["jwt"]
	return &config, nil
}

func readQuery(q backend.DataQuery) (*models.QueryModel, error) {
	queryModel := models.QueryModel{}
	if err := json.Unmarshal(q.JSON, &queryModel); err != nil {
		return nil, fmt.Errorf("failed to unmarshal query: %w", err)
	}
	return &queryModel, nil
}

// CheckHealth checks if the plugin is running properly
func (ds *GoogleSheetsDataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	res := &backend.CheckHealthResult{}

	// Just checking that the plugin exe is alive and running
	if req.PluginConfig.DataSourceConfig == nil {
		res.Status = backend.HealthStatusOk
		res.Message = "Plugin is Running"
		return res, nil
	}

	config, err := readConfig(req.PluginConfig)
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Invalid config"
		return res, nil
	}

	client, err := googlesheets.NewGoogleClient(ctx, config)
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Unable to create client"
		return res, nil
	}

	err = client.TestClient()
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Permissions check failed"
		return res, nil
	}

	res.Status = backend.HealthStatusOk
	res.Message = "Success"
	return res, nil
}

// QueryData queries for data.
func (ds *GoogleSheetsDataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	res := &backend.QueryDataResponse{}
	config, err := readConfig(req.PluginConfig)
	if err != nil {
		return nil, err
	}

	for _, q := range req.Queries {
		queryModel, err := readQuery(q)
		if err != nil {
			return nil, fmt.Errorf("failed to read query: %w", err)
		}

		if len(queryModel.Spreadsheet) < 1 {
			continue // not query really exists
		}

		frame, err := ds.googlesheet.Query(ctx, q.RefID, queryModel, config, q.TimeRange)
		if err != nil {
			backend.Logger.Error("Query failed", "refId", q.RefID, "error", err)
			return nil, err
		}

		res.Frames = append(res.Frames, []*data.Frame{frame}...)
	}

	return res, nil
}

func writeResult(rw http.ResponseWriter, path string, val interface{}, err error) {
	response := make(map[string]interface{})
	code := http.StatusOK
	if err != nil {
		response["error"] = err.Error()
		code = http.StatusBadRequest
	} else {
		response[path] = val
	}

	body, err := json.Marshal(response)
	if err != nil {
		body = []byte(err.Error())
		code = http.StatusInternalServerError
	}
	_, err = rw.Write(body)
	if err != nil {
		code = http.StatusInternalServerError
	}
	rw.WriteHeader(code)
}

func (ds *GoogleSheetsDataSource) handleResourceSpreadsheets(rw http.ResponseWriter, req *http.Request) {
	backend.Logger.Debug("Received resource call", "url", req.URL.String(), "method", req.Method)
	if req.Method != http.MethodGet {
		return
	}

	ctx := req.Context()
	config, err := readConfig(httpadapter.PluginConfigFromContext(ctx))
	if err != nil {
		writeResult(rw, "?", nil, err)
		return
	}

	res, err := ds.googlesheet.GetSpreadsheets(ctx, config)
	writeResult(rw, "spreadsheets", res, err)
}
