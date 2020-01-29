package googlesheets

type QueryModel struct {
	QueryType     string
	SpreadsheetID string
	Range         string
}

type GoogleSheetConfig struct {
	ApiKey   string `json:"apiKey"`
	AuthType string `json:"authType"`
	JwtFile  string `json:"jwtFile"`
}
