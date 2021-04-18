package main

import (
	"net/http"
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
)

func main() {
	mux := http.NewServeMux()
	ds := NewDataSource(mux)
	httpResourceHandler := httpadapter.New(mux)

	err := experimental.DoGRPC("google-sheets-datasource", datasource.ServeOpts{
		CallResourceHandler: httpResourceHandler,
		QueryDataHandler:    ds,
		CheckHealthHandler:  ds,
	})
	if err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
