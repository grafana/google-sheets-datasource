package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// QueryModel represents a spreadsheet query.
type QueryModel struct {
	Spreadsheet          string `json:"spreadsheet"`
	Range                string `json:"range"`
	CacheDurationSeconds int    `json:"cacheDurationSeconds"`
	UseTimeFilter        bool   `json:"useTimeFilter"`

	// Not from JSON
	TimeRange     backend.TimeRange `json:"-"`
	MaxDataPoints int64             `json:"-"`
}

// GetQueryModel returns the well typed query model
func GetQueryModel(query backend.DataQuery) (*QueryModel, error) {
	model := &QueryModel{}

	err := json.Unmarshal(query.JSON, &model)
	if err != nil {
		return nil, fmt.Errorf("error reading query: %s", err.Error())
	}

	// Copy directly from the well typed query
	model.TimeRange = query.TimeRange
	model.MaxDataPoints = query.MaxDataPoints
	return model, nil
}
