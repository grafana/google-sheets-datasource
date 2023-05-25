package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-google-sdk-go/pkg/utils"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// DatasourceSettings contains Google Sheets API authentication properties.
type DatasourceSettings struct {
	InstanceSettings   backend.DataSourceInstanceSettings
	AuthType           string `json:"authType"` // jwt | key | gce
	APIKey             string `json:"apiKey"`
	DefaultProject     string `json:"defaultProject"`
	JWT                string `json:"jwt"`
	ClientEmail        string `json:"clientEmail"`
	TokenURI           string `json:"tokenUri"`
	AuthenticationType string `json:"authenticationType"`
	PrivateKeyPath     string `json:"privateKeyPath"`

	// Saved in secure JSON
	PrivateKey string `json:"-"`
}

// LoadSettings gets the relevant settings from the plugin context
func LoadSettings(ctx backend.PluginContext) (*DatasourceSettings, error) {
	model := &DatasourceSettings{}

	settings := ctx.DataSourceInstanceSettings
	err := json.Unmarshal(settings.JSONData, &model)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}

	model.PrivateKey, err = utils.GetPrivateKey(settings)
	if err != nil {
		return model, err
	}

	model.APIKey = settings.DecryptedSecureJSONData["apiKey"]
	// Leaving this here for backward compatibility
	model.JWT = settings.DecryptedSecureJSONData["jwt"]
	model.InstanceSettings = *settings

	// Make sure that old settings are migrated to the new ones
	if model.AuthType != "" {
		model.AuthenticationType = model.AuthType
	}
	return model, nil
}
