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

// settingsExamples documents one worked configuration per authentication type,
// shaped as a datasource resource POST body: metadata.name, spec.jsonData, and
// secure.<field>.create for write-only secrets (privateKey, apiKey). See
// https://github.com/grafana/grafana/blob/1f6fc213e75436cf04af646735d21a0f33c27f5a/pkg/apis/datasource/v0alpha1/datasource.go#L17-L27.
var settingsExamples = &sdkSchema.SettingsExamples{
	Examples: map[string]*spec3.Example{
		"jwt": {
			ExampleProps: spec3.ExampleProps{
				Summary:     "Google JWT File",
				Description: "Service account credentials. Reads private and public spreadsheets. Set secure.privateKey.create to the `private_key` value from the service account JSON, or set spec.jsonData.privateKeyPath to a file on the Grafana server instead.",
				Value: map[string]any{
					"metadata": map[string]any{
						"name": "my-datasource",
					},
					"spec": map[string]any{
						"jsonData": map[string]any{
							"authenticationType": "jwt",
							"defaultProject":     "my-gcp-project-id",
							"clientEmail":        "grafana@my-gcp-project-id.iam.gserviceaccount.com",
							"tokenUri":           "https://oauth2.googleapis.com/token",
						},
					},
					"secure": map[string]any{
						"privateKey": map[string]any{
							"create": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
						},
					},
				},
			},
		},
		"key": {
			ExampleProps: spec3.ExampleProps{
				Summary:     "API Key",
				Description: "Simplest configuration, but only reads spreadsheets that are shared publicly.",
				Value: map[string]any{
					"metadata": map[string]any{
						"name": "my-datasource",
					},
					"spec": map[string]any{
						"jsonData": map[string]any{
							"authenticationType": "key",
						},
					},
					"secure": map[string]any{
						"apiKey": map[string]any{
							"create": "AIzaSy...",
						},
					},
				},
			},
		},
		"gce": {
			ExampleProps: spec3.ExampleProps{
				Summary:     "GCE Default Service Account",
				Description: "Credentials are retrieved from the GCE metadata server. Requires Grafana to be running on a Google Compute Engine virtual machine. No secrets are stored.",
				Value: map[string]any{
					"metadata": map[string]any{
						"name": "my-datasource",
					},
					"spec": map[string]any{
						"jsonData": map[string]any{
							"authenticationType": "gce",
							"defaultProject":     "my-gcp-project-id",
						},
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
