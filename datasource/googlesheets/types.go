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
	TimeColumn     valueOption `json:"timeColumn"`
	MetricColumn   valueOption `json:"metricColumn"`
}

type GoogleSheetConfig struct {
	ApiKey   string `json:"apiKey"`
	AuthType string `json:"authType"`
	JwtFile  string `json:"jwtFile"`
}
