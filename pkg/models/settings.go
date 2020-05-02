package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// GoogleSheetConfig contains Google Sheets API authentication properties.
type DatasourceSettings struct {
	AuthType string `json:"authType"` // jwt | key
	APIKey   string `json:"apiKey"`
	JWT      string `json:"jwt"`
}

func LoadSettings(ctx backend.PluginContext) (*DatasourceSettings, error) {
	model := &DatasourceSettings{}

	settings := ctx.DataSourceInstanceSettings
	err := json.Unmarshal(settings.JSONData, &model)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}

	model.APIKey = settings.DecryptedSecureJSONData["apiKey"]
	model.JWT = settings.DecryptedSecureJSONData["jwt"]

	return model, nil
}
