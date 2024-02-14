package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// QueryModel represents a spreadsheet query.
type QueryModel struct {
	// The google sheets spreadsheet ID
	Spreadsheet string `json:"spreadsheet"`

	// A1 notation
	Range string `json:"range,omitempty"`

	// Cache duration in seconds
	CacheDurationSeconds int `json:"cacheDurationSeconds,omitempty"`

	// Use the query time range to filer values from the table
	UseTimeFilter bool `json:"useTimeFilter,omitempty"`

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
