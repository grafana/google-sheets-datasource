package googlesheets

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
	ApiKey   string `json:"apiKey"`
	AuthType string `json:"authType"`
	JwtFile  string `json:"jwtFile"`
}
