package models

import (
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental/schema"
	"github.com/stretchr/testify/require"
)

func TestSchemaDefinitions(t *testing.T) {
	builder, err := schema.NewSchemaBuilder(
		schema.BuilderOptions{
			BasePackage: "github.com/grafana/google-sheets-datasource/pkg/models",
			CodePath:    "./",
		},
	)
	require.NoError(t, err)
	err = builder.AddQueries(schema.QueryTypeInfo{
		Name:   "default",
		GoType: reflect.TypeOf(&QueryModel{}),
		Examples: []schema.QueryExample{
			{
				Name: "public query",
				QueryPayload: QueryModel{
					Spreadsheet: "YourSheetID",
					Range:       "A1:D6",
				},
			},
		},
	})

	require.NoError(t, err)

	// Update the query schemas resource
	builder.UpdateQueryDefinition(t, "../../src/static/schema/query.schema.json")
}
