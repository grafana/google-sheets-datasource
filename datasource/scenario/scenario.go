package scenario

import (
	"context"
	"encoding/json"
	"fmt"

	ds "github.com/grafana/grafana-plugin-sdk-go/backend"
	df "github.com/grafana/grafana-plugin-sdk-go/dataframe"
	hclog "github.com/hashicorp/go-hclog"
)

// baseQuery includes the scenario type, and embeds sdk query
// object inside of it.
type baseQuery struct {
	Scenario  queryType `json:"scenario"`
	TimeRange ds.TimeRange
	logger    hclog.Logger
	ds.DataQuery
}

func (b *baseQuery) GetName() string {
	return string(b.Scenario)
}

func (b *baseQuery) SetLogger(logger hclog.Logger) {
	b.logger = logger
}

// queryType defines the expected type of scenario.
type queryType string

func (s *queryType) UnmarshalJSON(b []byte) error {
	var text string
	if err := json.Unmarshal(b, &text); err != nil {
		return err
	}
	switch queryType(text) {
	case typeCSVWaveQuery:
		*s = typeCSVWaveQuery
	case typeArrowFileQuery:
		*s = typeArrowFileQuery
	default:
		return fmt.Errorf("unknown scenario type '%v'", text)
	}
	return nil
}

// scenario is the interface for all scenarios.
type scenario interface {
	GetName() string
	Execute(ctx context.Context) ([]*df.Frame, error)
	SetLogger(logger hclog.Logger)
}

// UnmarshalQueries builds the appropriate query type objects based
// on the query's scenario property.
func UnmarshalQueries(logger hclog.Logger, queries []ds.DataQuery) ([]scenario, error) {
	procQueries := make([]scenario, len(queries))
	for idx := range queries {
		// baseQuery only used get queryType is and also set the time range
		// since it currently isn't per query in the JSON model (but is in the backend model).
		bq := baseQuery{TimeRange: queries[idx].TimeRange}
		err := json.Unmarshal(queries[idx].JSON, &bq)
		if err != nil {
			return nil, err
		}
		var s scenario
		switch bq.Scenario {
		case typeWalkQuery:
			q := &walkQuery{baseQuery: bq}
			err = json.Unmarshal(queries[idx].JSON, q)
			s = q
		case typeArrowFileQuery:
			q := &arrowFileQuery{baseQuery: bq}
			err = json.Unmarshal(queries[idx].JSON, q)
			s = q
		case typeCSVWaveQuery:
			q := &csvWaveQuery{baseQuery: bq}
			err = json.Unmarshal(queries[idx].JSON, q)
			s = q
		default:
			return nil, fmt.Errorf("unsupported scenario type %v", bq.Scenario)
		}
		if err != nil {
			return nil, err
		}
		s.SetLogger(logger)
		procQueries[idx] = s

	}
	return procQueries, nil
}
