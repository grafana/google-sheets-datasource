package models

import (
	"reflect"
	"testing"

	data "github.com/grafana/grafana-plugin-sdk-go/experimental/apis/data/v0alpha1"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/schemabuilder"
	"github.com/stretchr/testify/require"
)

func TestSchemaDefinitions(t *testing.T) {
	builder, err := schemabuilder.NewSchemaBuilder(
		schemabuilder.BuilderOptions{
			PluginID: []string{"grafana-googlesheets-datasource"},
			ScanCode: []schemabuilder.CodePaths{{
				BasePackage: "github.com/grafana/google-sheets-datasource/pkg/models",
				CodePath:    "./",
			}},
		},
	)
	require.NoError(t, err)
	err = builder.AddQueries(schemabuilder.QueryTypeInfo{
		Name:   "default",
		GoType: reflect.TypeOf(&QueryModel{}),
		Examples: []data.QueryExample{
			{
				Name: "public query",
				SaveModel: data.AsUnstructured(QueryModel{
					Spreadsheet: "YourSheetID",
					Range:       "A1:D6",
				}),
			},
		},
	})

	require.NoError(t, err)

	// Update the query schemas resource
	builder.UpdateQueryDefinition(t, "../../src/static/schema")
}
