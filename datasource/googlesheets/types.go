package googlesheets

type QueryModel struct {
	QueryType     string
	SpreadsheetID string
	Range         string
}

type GoogleSheetConfig struct {
	ApiKey string `json:"apiKey"`
}
