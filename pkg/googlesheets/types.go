package googlesheets

import "encoding/json"

type QueryModel struct {
	QueryType   string `json:"queryType"`
	Spreadsheet struct {
		ID   string `json:"value"`
		Name string `json:"label"`
	} `json:"Spreadsheet"`
	Range                string `json:"range"`
	CacheDurationSeconds int    `json:"cacheDurationSeconds"`
}

type GoogleSheetConfig struct {
	ApiKey   string          `json:"apiKey"`
	AuthType string          `json:"authType"`
	JWT      json.RawMessage `json:"jwt"`
}

// type googleClient interface {
// 	createClient
// }
