package googlesheets

import (
	"encoding/json"

	"google.golang.org/api/sheets/v4"
)

type Spreadsheet struct {
	ID   string `json:"value"`
	Name string `json:"label"`
}

type QueryModel struct {
	QueryType            string      `json:"queryType"`
	Spreadsheet          Spreadsheet `json:"Spreadsheet"`
	Range                string      `json:"range"`
	CacheDurationSeconds int         `json:"cacheDurationSeconds"`
}

type GoogleSheetConfig struct {
	ApiKey   string          `json:"apiKey"`
	AuthType string          `json:"authType"`
	JWT      json.RawMessage `json:"jwt"`
}

type client interface {
	GetSpreadsheet(spreadSheetID string, sheetRange string, includeGridData bool) (*sheets.Spreadsheet, error)
}
