package main

import (
	"encoding/json"
	"fmt"
	"os"

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

func (gsd *GoogleSheetsDataSource) DataQuery(ctx context.Context, req *backend.DataQueryRequest) (*backend.DataQueryResponse, error) {
	res := &backend.DataQueryResponse{}
	for _, q := range req.Queries {
		queryModel := &gs.QueryModel{}
		err := json.Unmarshal(q.JSON, &queryModel)
		if err != nil {
			gsd.logger.Error("Failed to unmarshal query: %v", err.Error())
		}

		var frame *df.Frame
		switch queryModel.QueryType {
		case "testAPI":
			frame, err = gs.TestAPI()
		case "query":
			apiKey, _ := os.LookupEnv(variableName)
			frame, err = gs.Query(ctx, q.RefID, queryModel, &gs.GoogleSheetConfig{ApiKey: apiKey})
		default:
			return nil, fmt.Errorf("Invalid query type")
		}

		if err != nil {
			gsd.logger.Error("Failed to execute query: %v", err.Error())
		} else {
			res.Frames = append(res.Frames, frame)
		}

	}

	return res, nil
}
