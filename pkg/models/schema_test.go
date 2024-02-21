package models

import (
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental/spec"
	"github.com/stretchr/testify/require"
)

func TestSchemaDefinitions(t *testing.T) {
	builder, err := spec.NewSchemaBuilder(
		spec.BuilderOptions{
			BasePackage: "github.com/grafana/google-sheets-datasource/pkg/models",
			CodePath:    "./",
		},
	)
	require.NoError(t, err)
	err = builder.AddQueries(spec.QueryTypeInfo{
		Name:   "default",
		GoType: reflect.TypeOf(&QueryModel{}),
		Examples: []spec.QueryExample{
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
	builder.UpdateQueryDefinition(t, "../../src/static/schema")
}
