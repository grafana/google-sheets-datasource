package core

// QueryModel represents a spreadsheet query.
// This needs to match the request sent from:
// https://github.com/grafana/google-sheets-datasource/blob/master/src/types.ts#L46
type QueryModel struct {
	Spreadsheet          string `json:"spreadsheet"`
	Range                string `json:"range"`
	CacheDurationSeconds int    `json:"cacheDurationSeconds"`
	UseTimeFilter        bool   `json:"useTimeFilter"`
}

// GoogleSheetConfig contains Google Sheets API authentication properties.
type GoogleSheetConfig struct {
	AuthType string `json:"authType"` // jwt | key
	APIKey   string `json:"apiKey"`   // Saved in secure JSON fields
	JWT      string `json:"jwt"`      // Saved in secure JSON fields
}
