package main

import (
	"os"

	"github.com/grafana/google-sheets-datasource/pkg/ext"
	"github.com/grafana/google-sheets-datasource/pkg/googlesheets"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

func main() {
	// check if k8s is enabled
	for _, v := range os.Args {
		if v == "k8s" {
			stopCh := genericapiserver.SetupSignalHandler()
			options := ext.NewPluginAggregatedServerOptions(os.Stdout, os.Stderr)

			if err := options.Complete(); err != nil {
				panic(err)
			}

			if err := options.Run(stopCh); err != nil {
				panic(err)
			}
			os.Exit(1)
		}
	}

	// The standard startup path
	if err := datasource.Manage("google-sheets-datasource", googlesheets.NewDatasource, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
