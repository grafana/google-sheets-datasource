package googlesheets

import (
	"google.golang.org/api/sheets/v4"
)

// Spreadsheet represents a Google spreadsheet.
type Spreadsheet struct {
	ID   string `json:"value"`
	Name string `json:"label"`
}

// QueryModel represents a spreadsheet query.
type QueryModel struct {
	Spreadsheet          Spreadsheet `json:"Spreadsheet"`
	Range                string      `json:"range"`
	CacheDurationSeconds int         `json:"cacheDurationSeconds"`
}

// GoogleSheetConfig contains Google Sheets API authentication properties.
type GoogleSheetConfig struct {
	APIKey   string `json:"apiKey"`
	AuthType string `json:"authType"`
	JWT      string `json:"jwt"`
}

type client interface {
	GetSpreadsheet(spreadSheetID string, sheetRange string, includeGridData bool) (*sheets.Spreadsheet, error)
}
