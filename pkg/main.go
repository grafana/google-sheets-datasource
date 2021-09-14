package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
)

func main() {
	// mux := http.NewServeMux()
	// httpResourceHandler := httpadapter.New(mux)

	if err := datasource.Manage("google-sheets-datasource", NewDataSource, datasource.ManageOpts{}); err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
