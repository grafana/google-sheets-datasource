package models

import (
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental/resource"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/resource/schemabuilder"
	"github.com/stretchr/testify/require"
)

func TestSchemaDefinitions(t *testing.T) {
	builder, err := schemabuilder.NewSchemaBuilder(
		schemabuilder.BuilderOptions{
			PluginID:    []string{"grafana-googlesheets-datasource"},
			BasePackage: "github.com/grafana/google-sheets-datasource/pkg/models",
			CodePath:    "./",
		},
	)
	require.NoError(t, err)
	err = builder.AddQueries(schemabuilder.QueryTypeInfo{
		Name:   "default",
		GoType: reflect.TypeOf(&QueryModel{}),
		Examples: []resource.QueryExample{
			{
				Name: "public query",
				SaveModel: QueryModel{
					Spreadsheet: "YourSheetID",
					Range:       "A1:D6",
				},
			},
		},
	})

	require.NoError(t, err)

	// Update the query schemas resource
	builder.UpdateQueryDefinition(t, "../../src/static/schema")
}
