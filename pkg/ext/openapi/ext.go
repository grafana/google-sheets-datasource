package openapi

import (
	generatedopenapi "github.com/grafana/google-sheets-datasource/pkg/client/openapi"
	"k8s.io/kube-openapi/pkg/common"
	spec "k8s.io/kube-openapi/pkg/validation/spec"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets/v1.DatasourceQuery": schema_pkg_apis_googlesheets_v1_DatasourceQuery(ref),
	}
}

func schema_pkg_apis_googlesheets_v1_DatasourceQuery(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "DatasourceQuery defines the observed state of Datasource",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"Results": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"map[string]backend.DataResponse"},
							Format:  "",
							Ref:     ref("github.com/grafana/grafana-plugin-sdk-go/backend"),
						},
					},
				},
			},
		},
	}
}

func AttachExtDefinitionsToGenerated(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	extended := make(map[string]common.OpenAPIDefinition)
	for k, v := range generatedopenapi.GetOpenAPIDefinitions(ref) {
		extended[k] = v
	}

	for k, v := range GetOpenAPIDefinitions(ref) {
		extended[k] = v
	}
	return extended

}
