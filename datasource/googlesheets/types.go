package googlesheets

type valueOption struct {
	Label string
	Value int
}

type QueryModel struct {
	ResultFormat   string
	QueryType      string
	SpreadsheetID  string
	Range          string
	MajorDimension string
	TimeColumn     valueOption   `json:"timeColumn"`
	MetricColumns  []valueOption `json:"metricColumns"`
}

type GoogleSheetConfig struct {
	ApiKey   string `json:"apiKey"`
	AuthType string `json:"authType"`
	JwtFile  string `json:"jwtFile"`
}
