package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	gc "github.com/grafana/google-sheets-datasource/pkg/googlesheets/googleclient"

	"github.com/grafana/google-sheets-datasource/pkg/googlesheets"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpresource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	p "github.com/grafana/grafana-plugin-sdk-go/backend/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"

	"context"
)

const metricNamespace = "sheets_datasource"

func main() {
	// Setup the plugin environment
	_ = p.SetupPluginEnvironment("google-sheets-datasource")
	pluginLogger := log.New()

	mux := http.NewServeMux()
	ds := Init(pluginLogger, mux)
	httpResourceHandler := httpresource.New(mux)

	err := backend.Serve(backend.ServeOpts{
		CallResourceHandler: httpResourceHandler,
		QueryDataHandler:    ds,
		CheckHealthHandler:  ds,
	})
	if err != nil {
		pluginLogger.Error(err.Error())
		os.Exit(1)
	}
}

// Init creates the google sheets datasource and sets up all the routes
func Init(logger log.Logger, mux *http.ServeMux) *GoogleSheetsDataSource {
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
		logger: logger,
		googlesheet: &googlesheets.GoogleSheets{
			Cache:  cache,
			Logger: logger,
		},
	}

	mux.HandleFunc("/spreadsheets", ds.handleResourceSpreadsheets)
	return ds
}

// GoogleSheetsDataSource handler for google sheets
type GoogleSheetsDataSource struct {
	logger      log.Logger
	googlesheet *googlesheets.GoogleSheets
}

func getConfig(pluginConfig backend.PluginConfig) (*googlesheets.GoogleSheetConfig, error) {
	config := googlesheets.GoogleSheetConfig{}
	if err := json.Unmarshal(pluginConfig.DataSourceConfig.JSONData, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal DataSourceInfo json: %w", err)
	}
	config.APIKey = pluginConfig.DataSourceConfig.DecryptedSecureJSONData["apiKey"]
	config.JWT = pluginConfig.DataSourceConfig.DecryptedSecureJSONData["jwt"]
	return &config, nil
}

// CheckHealth checks if the plugin is running properly
func (plugin *GoogleSheetsDataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	res := &backend.CheckHealthResult{}

	if true {
		res.Status = backend.HealthStatusError
		res.Message = "Plugin is Running"
		return res, nil
	}

	// Just checking that the plugin exe is alive and running
	if req.PluginConfig.DataSourceConfig == nil {
		res.Status = backend.HealthStatusOk
		res.Message = "Plugin is Running"
		return res, nil
	}

	config, err := getConfig(req.PluginConfig)
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Invalid config"
		return res, nil
	}

	client, err := gc.New(ctx, gc.NewAuth(config.APIKey, config.AuthType, config.JWT))
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
func (plugin *GoogleSheetsDataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	res := &backend.QueryDataResponse{}
	config, err := getConfig(req.PluginConfig)
	if err != nil {
		return nil, err
	}

	for _, q := range req.Queries {
		queryModel := &googlesheets.QueryModel{}
		if err := json.Unmarshal(q.JSON, &queryModel); err != nil {
			return nil, fmt.Errorf("failed to unmarshal query: %w", err)
		}

		if len(queryModel.Spreadsheet) < 1 {
			continue // not query really exists
		}

		frame, err := plugin.googlesheet.Query(ctx, q.RefID, queryModel, config, q.TimeRange)
		if err != nil {
			plugin.logger.Error("Query failed", "refId", q.RefID, "error", err)
			// TEMP: at the moment, the only way to return an error is by using meta
			res.Metadata = map[string]string{"error": err.Error()}
			continue
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

func (plugin *GoogleSheetsDataSource) handleResourceSpreadsheets(rw http.ResponseWriter, req *http.Request) {
	plugin.logger.Debug("Received resource call", "url", req.URL.String(), "method", req.Method)
	if req.Method != http.MethodGet {
		return
	}

	ctx := req.Context()
	config, err := getConfig(httpresource.PluginConfigFromContext(ctx))
	if err != nil {
		writeResult(rw, "?", nil, err)
		return
	}

	res, err := plugin.googlesheet.GetSpreadsheets(ctx, config)
	writeResult(rw, "spreadsheets", res, err)
}
