package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/grafana/google-sheets-datasource/pkg/googlesheets"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpresource"
	df "github.com/grafana/grafana-plugin-sdk-go/dataframe"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"

	hclog "github.com/hashicorp/go-hclog"

	"context"
)

const metricNamespace = "sheets_datasource"

// Init creates the google sheets datasource and sets up all the routes
func Init(logger hclog.Logger, mux *http.ServeMux) *GoogleSheetsDataSource {
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

	mux.HandleFunc("/test", ds.handleResourceTest)
	mux.HandleFunc("/spreadsheets", ds.handleResourceSpreadsheets)
	return ds
}

// GoogleSheetsDataSource handler for google sheets
type GoogleSheetsDataSource struct {
	logger      hclog.Logger
	googlesheet *googlesheets.GoogleSheets
}

func getConfig(pluginConfig backend.PluginConfig) (*googlesheets.GoogleSheetConfig, error) {
	config := googlesheets.GoogleSheetConfig{}
	if err := json.Unmarshal(pluginConfig.JSONData, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal DataSourceInfo json: %w", err)
	}

	config.APIKey = pluginConfig.DecryptedSecureJSONData["apiKey"]
	config.JWT = pluginConfig.DecryptedSecureJSONData["jwt"]
	return &config, nil
}

// DataQuery queries for data.
func (ds *GoogleSheetsDataSource) DataQuery(ctx context.Context, req *backend.DataQueryRequest) (*backend.DataQueryResponse, error) {
	res := &backend.DataQueryResponse{}
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

		frame, err := ds.googlesheet.Query(ctx, q.RefID, queryModel, config, q.TimeRange)
		if err != nil {
			ds.logger.Error("Query failed", "refId", q.RefID, "error", err)
			// TEMP: at the moment, the only way to return an error is by using meta
			res.Metadata = map[string]string{"error": err.Error()}
			continue
		}

		res.Frames = append(res.Frames, []*df.Frame{frame}...)
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
	rw.Write(body)
	rw.WriteHeader(code)
}

func (ds *GoogleSheetsDataSource) handleResourceSpreadsheets(rw http.ResponseWriter, req *http.Request) {
	ds.logger.Debug("Received resource call", "url", req.URL.String(), "method", req.Method)
	if req.Method != http.MethodGet {
		return
	}

	ctx := req.Context()
	config, err := getConfig(httpresource.PluginConfigFromContext(ctx))
	if err != nil {
		writeResult(rw, "?", nil, err)
		return
	}

	res, err := ds.googlesheet.GetSpreadsheets(ctx, config)
	writeResult(rw, "spreadsheets", res, err)
}

func (ds *GoogleSheetsDataSource) handleResourceTest(rw http.ResponseWriter, req *http.Request) {
	ds.logger.Debug("Received resource call", "url", req.URL.String(), "method", req.Method)
	if req.Method != http.MethodGet {
		return
	}

	ctx := req.Context()
	config, err := getConfig(httpresource.PluginConfigFromContext(ctx))
	if err != nil {
		writeResult(rw, "?", nil, err)
		return
	}

	err = ds.googlesheet.TestAPI(ctx, config)
	writeResult(rw, "test", "OK", err)
}
