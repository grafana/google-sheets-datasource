package main

import (
	"context"
	"github.com/grafana/google-sheets-datasource/pkg/ext"
	"os"

	"github.com/grafana/google-sheets-datasource/pkg/googlesheets"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func main() {
	err := ext.Start(context.TODO())
	if err != nil {
		panic(err)
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
