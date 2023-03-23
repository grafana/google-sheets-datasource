package models

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
	TokenUri           string `json:"tokenUri"`
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

	// Check if a private key path was provided. Fall back to the plugin's default method
	// of an inline private key
	if model.PrivateKeyPath != "" {
		privateKey, err := readPrivateKeyFromFile(model.PrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("could not write private key to DataSourceInfo json: %w", err)
		}

		model.PrivateKey = privateKey
	} else {
		privateKey := settings.DecryptedSecureJSONData["privateKey"]
		// React might escape newline characters like this \\n so we need to handle that
		model.PrivateKey = strings.ReplaceAll(privateKey, "\\n", "\n")
	}

	model.APIKey = settings.DecryptedSecureJSONData["apiKey"]
	// Leaving this here for backward compatibility
	model.JWT = settings.DecryptedSecureJSONData["jwt"]
	model.InstanceSettings = *settings
	model.AuthenticationType = model.AuthType

	return model, nil
}

func readPrivateKeyFromFile(rsaPrivateKeyLocation string) (string, error) {
	if rsaPrivateKeyLocation == "" {
		return "", fmt.Errorf("missing file location for private key")
	}

	privateKey, err := os.ReadFile(rsaPrivateKeyLocation)
	if err != nil {
		return "", fmt.Errorf("could not read private key file from file system: %w", err)
	}

	return string(privateKey), nil
}
