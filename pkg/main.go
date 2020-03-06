package main

import (
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpresource"
	hclog "github.com/hashicorp/go-hclog"
)

func main() {
	// Setup the plugin environment
	pluginLogger := setupPluginEnvironment("google-sheets-datasource")

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

// TODO: move this function to the SDK
func setupPluginEnvironment(pluginID string) hclog.Logger {
	pluginLogger := hclog.New(&hclog.LoggerOptions{
		Name: pluginID,
		// TODO: How to make level configurable?
		Level:      hclog.LevelFromString("DEBUG"),
		JSONFormat: true,
		Color:      hclog.ColorOff,
	})

	// Enable profiler
	profilerEnabled := false
	if value, ok := os.LookupEnv("GF_PLUGINS_PROFILER"); ok {
		// compare value to plugin name
		if value == pluginID {
			profilerEnabled = true
		}
	}
	pluginLogger.Info("Profiler", "enabled", profilerEnabled)
	if profilerEnabled {
		r := http.NewServeMux()
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)

		go func() {
			if err := http.ListenAndServe(":6060", r); err != nil {
				pluginLogger.Error("Error Running profiler: %s", err.Error())
			}
		}()
	}
	return pluginLogger
}
