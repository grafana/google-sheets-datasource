package main

import (
	"os"

	"github.com/grafana/google-sheets-datasource/pkg/ext"
	genericapiserver "k8s.io/apiserver/pkg/server"
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
}
