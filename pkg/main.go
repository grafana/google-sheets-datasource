package main

import (
	"os"

	"github.com/grafana/google-sheets-datasource/pkg/ext"
	"github.com/grafana/google-sheets-datasource/pkg/googlesheets"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func main() {
	// check if k8s is enabled
	for _, v := range os.Args {
		if v == "k8s" {
			err := ext.RunServer()
			if err != nil {
				panic(err)
			}
			os.Exit(1)
		}
	}

	if err := datasource.Manage("google-sheets-datasource", googlesheets.NewDatasource, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}

//
// Group:
//   google-sheets.ext.grafana.com
// Kind:
//   datasource ??
//
