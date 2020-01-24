package main

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/test-datasource/datasource/datasource/scenario"

	hclog "github.com/hashicorp/go-hclog"
)

const pluginID = "backend"

var pluginLogger = hclog.New(&hclog.LoggerOptions{
	Name:  pluginID,
	Level: hclog.LevelFromString("DEBUG"),
})

const metricNamespace = "test_datasource"

var (
	queriesTotal *prometheus.CounterVec
)

func main() {
	queriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "data_query_total",
			Help:      "data query counter",
			Namespace: metricNamespace,
		},
		[]string{"scenario"},
	)
	prometheus.MustRegister(queriesTotal)

	ds := &testDatasource{
		logger: pluginLogger,
	}
	err := backend.Serve(backend.ServeOpts{
		DataQueryHandler: ds,
	})
	if err != nil {
		pluginLogger.Error(err.Error())
	}
}

type testDatasource struct {
	logger hclog.Logger
}

func (td *testDatasource) DataQuery(ctx context.Context, req *backend.DataQueryRequest) (*backend.DataQueryResponse, error) {
	scenarioQueries, err := scenario.UnmarshalQueries(td.logger, req.Queries)
	if err != nil {
		return nil, err
	}

	res := backend.DataQueryResponse{}
	for _, q := range scenarioQueries {
		queriesTotal.WithLabelValues(q.GetName()).Inc()
		frames, err := q.Execute(ctx)
		if err != nil {
			return nil, err
		}
		res.Frames = append(res.Frames, frames...)
	}
	return &res, nil
}
