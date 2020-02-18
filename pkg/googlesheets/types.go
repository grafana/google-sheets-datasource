package googlesheets

type QueryModel struct {
	QueryType            string
	SpreadsheetID        string
	Range                string
	CacheDurationSeconds int
}

type GoogleSheetConfig struct {
	ApiKey   string `json:"apiKey"`
	AuthType string `json:"authType"`
	JwtFile  string `json:"jwtFile"`
}
