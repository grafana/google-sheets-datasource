package main

import (
	"net/http"
	"os"

	"github.com/grafana/google-sheets-datasource/pkg/datasource"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
)

func main() {
	// Setup the plugin environment
	_ = backend.SetupPluginEnvironment("google-sheets-datasource")
	pluginLogger := log.New()

	mux := http.NewServeMux()
	ds := datasource.Init(pluginLogger, mux)
	httpResourceHandler := httpadapter.New(mux)

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
