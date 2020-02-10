package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	gs "github.com/grafana/google-sheets-datasource/datasource/googlesheets"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	df "github.com/grafana/grafana-plugin-sdk-go/dataframe"

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
	ds := &GoogleSheetsDataSource{
		logger: pluginLogger,
	}
	err := backend.Serve(backend.ServeOpts{
		CallResourceHandler: ds,
		DataQueryHandler:    ds,
	})
	if err != nil {
		pluginLogger.Error(err.Error())
	}
}

type GoogleSheetsDataSource struct {
	plugin.NetRPCUnsupportedPlugin
	logger hclog.Logger
}

func (gsd *GoogleSheetsDataSource) DataQuery(ctx context.Context, req *backend.DataQueryRequest) (*backend.DataQueryResponse, error) {
	gsd.logger.Debug("GoogleSheetsDatasource", "DataQuery")
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
			frames, err = gs.TestAPI(ctx, &config)
		case "query":
			frames, err = gs.Query(ctx, q.RefID, queryModel, &config, q.TimeRange, gsd.logger)
		case "getSpreadsheets":
			a, _ := gs.GetSpreadsheets(ctx, &config, gsd.logger)
			frame := df.New("getHeader")
			res.Metadata = a
			frames = []*df.Frame{frame}
		case "getHeaders":
			// res.Metadata = make(map[string]string)
			// frame = df.New("getHeader")
			// frame.RefID = q.RefID
			// frame.Meta = &df.QueryResultMeta{Custom: make(map[string]interface{})}
			// res, err := gs.GetHeaders(ctx, queryModel, &config, gsd.logger)
			// if err == nil {
			// 	frame.Meta.Custom["headers"] = res
			// }
			frame := df.New("getHeaders")
			headers, _ := gs.GetHeaders(ctx, queryModel, &config, gsd.logger)
			res.Metadata = make(map[string]string)
			for i, header := range headers {
				res.Metadata[fmt.Sprint(i)] = header
			}
			frames = []*df.Frame{frame}
		default:
			return nil, fmt.Errorf("Invalid query type")
		}

		if err != nil {
			gsd.logger.Debug("metric error=", err)
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

func (gsd *GoogleSheetsDataSource) CallResource(ctx context.Context, req *backend.CallResourceRequest) (*backend.CallResourceResponse, error) {
	gsd.logger.Debug("aaaaa7: ")
	gsd.logger.Debug(spew.Sdump(req.PluginConfig))
	response := make(map[string]interface{})

	config := gs.GoogleSheetConfig{}
	err := json.Unmarshal(req.PluginConfig.JSONData, &config)
	if err != nil {
		gsd.logger.Error("Could not unmarshal DataSourceInfo json", "Error", err)
		return nil, err
	}

	var metaQuery = &gs.MetaQuery{}
	if len(req.Body) > 0 {
		err := json.Unmarshal(req.Body, &metaQuery)
		if err != nil {
			return nil, err
		}
	}

	switch metaQuery.QueryType {
	case "getHeaders":
		headers, err := gs.GetHeaders(ctx, &metaQuery.Query, &config, gsd.logger)
		if err != nil {
			gsd.logger.Error("Failed to get headers: %v", err.Error())
			return nil, fmt.Errorf("Invalid query")
		}

		response["headers"] = headers
	}

	body, err := json.Marshal(&response)
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
