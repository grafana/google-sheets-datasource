package googlesheets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/grafana/google-sheets-datasource/pkg/models"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/patrickmn/go-cache"
)

var (
	_ backend.QueryDataHandler    = (*Datasource)(nil)
	_ backend.CheckHealthHandler  = (*Datasource)(nil)
	_ backend.CallResourceHandler = (*Datasource)(nil)
)

type Datasource struct {
	googlesheets *GoogleSheets

	backend.CallResourceHandler
}

// NewDatasource creates a new Google Sheets datasource instance.
func NewDatasource(_ context.Context, _ backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	ds := &Datasource{
		googlesheets: &GoogleSheets{Cache: cache.New(300*time.Second, 5*time.Second)},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/spreadsheets", ds.handleResourceSpreadsheets)
	ds.CallResourceHandler = httpadapter.New(mux)

	return ds, nil
}

// CheckHealth checks if the datasource is working.
func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	logger := backend.Logger.FromContext(ctx)
	res := &backend.CheckHealthResult{}
	logger.Debug("CheckHealth called")
	config, err := models.LoadSettings(req.PluginContext)

	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Unable to load settings"
		logger.Debug(err.Error())
		return res, nil
	}

	client, err := NewGoogleClient(ctx, *config)
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Unable to create client"
		logger.Debug(err.Error())
		return res, nil
	}

	err = client.TestClient()
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Permissions check failed"
		logger.Debug(err.Error())
		return res, nil
	}

	res.Status = backend.HealthStatusOk
	res.Message = "Success"
	return res, nil
}

// QueryData handles queries to the datasource.
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	logger := backend.Logger.FromContext(ctx)
	// create response struct
	response := backend.NewQueryDataResponse()

	logger.Debug("QueryData called", "numQueries", len(req.Queries))

	config, err := models.LoadSettings(req.PluginContext)
	if err != nil {
		return nil, err
	}

	for _, q := range req.Queries {
		queryModel, err := models.GetQueryModel(q)
		if err != nil {
			return nil, fmt.Errorf("failed to read query: %w", err)
		}

		if len(queryModel.Spreadsheet) < 1 {
			continue // not query really exists
		}
		dr := d.googlesheets.Query(ctx, q.RefID, queryModel, *config, q.TimeRange)
		if dr.Error != nil {
			if dr.ErrorSource == backend.ErrorSourceDownstream {
				// For downstream errors, we log them as warnings as they are not caused by the plugin itself
				logger.Debug("Query failed", "refId", q.RefID, "error", dr.Error, "errorsource", dr.ErrorSource)
			} else {
				logger.Error("Query failed", "refId", q.RefID, "error", dr.Error, "errorsource", dr.ErrorSource)
			}
		}
		response.Responses[q.RefID] = dr
	}

	return response, nil
}

func writeResult(rw http.ResponseWriter, path string, val any, err error) {
	response := make(map[string]any)
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

func (d *Datasource) handleResourceSpreadsheets(rw http.ResponseWriter, req *http.Request) {
	log.DefaultLogger.Debug("Received resource call", "url", req.URL.String())
	if req.Method != http.MethodGet {
		return
	}

	ctx := req.Context()
	config, err := models.LoadSettings(backend.PluginConfigFromContext(ctx))
	if err != nil {
		writeResult(rw, "?", nil, err)
		return
	}

	res, err := d.googlesheets.GetSpreadsheets(ctx, *config)
	writeResult(rw, "spreadsheets", res, err)
}
