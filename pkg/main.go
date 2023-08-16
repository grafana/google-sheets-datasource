package main

import (
	"github.com/grafana/google-sheets-datasource/pkg/ext"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"os"

	"github.com/grafana/google-sheets-datasource/pkg/googlesheets"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func main() {
	stopCh := genericapiserver.SetupSignalHandler()
	options := ext.NewPluginAggregatedServerOptions(os.Stdout, os.Stderr)

	if err := options.Complete(); err != nil {
		panic(err)
	}

	if err := options.Run(stopCh); err != nil {
		panic(err)
	}

	if err := datasource.Manage("google-sheets-datasource", googlesheets.NewDatasource, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}

//
// Group:
//   googlesheets.ext.grafana.com
// Kind:
//   datasource ??
//
