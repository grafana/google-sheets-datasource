package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// Fork to temporarily use private key insensitively for POC from SDK util

type JSONData struct {
	PrivateKey string `json:"privateKey"`
}

func getPrivateKey(settings *backend.DataSourceInstanceSettings) (string, error) {
	jsonData := JSONData{}

	if err := json.Unmarshal(settings.JSONData, &jsonData); err != nil {
		return "", fmt.Errorf("could not unmarshal DataSourceInfo json: %w", err)
	}

	if jsonData.PrivateKey != "" {
		// React might escape newline characters like this \\n so we need to handle that
		return strings.ReplaceAll(jsonData.PrivateKey, "\\n", "\n"), nil
	}

	return "", fmt.Errorf("Missing private key in json data")
}
