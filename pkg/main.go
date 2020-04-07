package main

import (
	"net/http"
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
)

func main() {
	backend.SetupPluginEnvironment("google-sheets-datasource")

	mux := http.NewServeMux()
	ds := NewDataSource(mux)
	httpResourceHandler := httpadapter.New(mux)

	err := backend.Serve(backend.ServeOpts{
		CallResourceHandler: httpResourceHandler,
		QueryDataHandler:    ds,
		CheckHealthHandler:  ds,
	})
	if err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
