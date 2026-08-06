package schema_test

import (
	_ "embed"
	"testing"

	"github.com/grafana/dsconfig/schema"
	"github.com/grafana/google-sheets-datasource/pkg/models"
	sdkSchema "github.com/grafana/grafana-plugin-sdk-go/experimental/pluginschema"
	"k8s.io/kube-openapi/pkg/spec3"
)

//go:embed dsconfig.json
var configSchemaJSON []byte

// settingsExamples documents one worked configuration per authentication type.
// Only jsonData is shown; secrets (privateKey, apiKey) are written to
// secureJsonData and are never readable back.
var settingsExamples = &sdkSchema.SettingsExamples{
	Examples: map[string]*spec3.Example{
		"jwt": {
			ExampleProps: spec3.ExampleProps{
				Summary:     "Google JWT File",
				Description: "Service account credentials. Reads private and public spreadsheets. Set secureJsonData.privateKey to the `private_key` value from the service account JSON, or set jsonData.privateKeyPath to a file on the Grafana server instead.",
				Value: map[string]any{
					"jsonData": map[string]any{
						"authenticationType": "jwt",
						"defaultProject":     "my-gcp-project-id",
						"clientEmail":        "grafana@my-gcp-project-id.iam.gserviceaccount.com",
						"tokenUri":           "https://oauth2.googleapis.com/token",
					},
					"secureJsonData": map[string]any{
						"privateKey": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
					},
				},
			},
		},
		"key": {
			ExampleProps: spec3.ExampleProps{
				Summary:     "API Key",
				Description: "Simplest configuration, but only reads spreadsheets that are shared publicly.",
				Value: map[string]any{
					"jsonData": map[string]any{
						"authenticationType": "key",
					},
					"secureJsonData": map[string]any{
						"apiKey": "AIzaSy...",
					},
				},
			},
		},
		"gce": {
			ExampleProps: spec3.ExampleProps{
				Summary:     "GCE Default Service Account",
				Description: "Credentials are retrieved from the GCE metadata server. Requires Grafana to be running on a Google Compute Engine virtual machine. No secrets are stored.",
				Value: map[string]any{
					"jsonData": map[string]any{
						"authenticationType": "gce",
						"defaultProject":     "my-gcp-project-id",
					},
				},
			},
		},
	},
}

//go:generate go test -run TestPlugin -generateArtifacts
func TestPlugin(t *testing.T) {
	schema.RunPluginTests(t, schema.PluginUnderTest{
		ID:                "grafana-googlesheets-datasource",
		ConfigSchemaJSON:  configSchemaJSON,
		SettingsJSONModel: models.DatasourceSettings{},
		SecureKeys:        []string{"privateKey", "apiKey", "jwt"},
		SettingsExamples:  settingsExamples,
	})
}
