package main

import (
	"net/http"
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpresource"
	hclog "github.com/hashicorp/go-hclog"
)

const (
	pluginID = "google-sheets-datasource"
)

func main() {
	pluginLogger := hclog.New(&hclog.LoggerOptions{
		Name: pluginID,
		// TODO: How to make level configurable?
		Level:      hclog.LevelFromString("DEBUG"),
		JSONFormat: true,
		Color:      hclog.ColorOff,
	})

	mux := http.NewServeMux()
	ds := Init(pluginLogger, mux)
	httpResourceHandler := httpresource.New(mux)

	err := backend.Serve(backend.ServeOpts{
		CallResourceHandler: httpResourceHandler,
		DataQueryHandler:    ds,
	})
	if err != nil {
		pluginLogger.Error(err.Error())
		os.Exit(1)
	}
}
