package googlesheets

type GoogleSheetRangeInfo struct {
	QueryType     string
	SpreadsheetID string
	Range         string
}

type GoogleSheetConfig struct {
	ApiKey string `json:"path"`
}
