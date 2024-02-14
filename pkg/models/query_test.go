package models

import (
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental/query"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/query/schema"
	"github.com/stretchr/testify/require"
)

func TestQueryTypeDefinitions(t *testing.T) {
	builder, err := schema.NewBuilder(
		schema.BuilderOptions{
			BasePackage: "github.com/grafana/google-sheets-datasource/pkg/models",
			CodePath:    "./",
		},
		schema.QueryTypeInfo{
			Name:   "default",
			GoType: reflect.TypeOf(&QueryModel{}),
			Examples: []query.QueryExample{
				{
					Name: "public query",
					Query: QueryModel{
						Spreadsheet: "YourSheetID",
						Range:       "A1:D6",
					},
				},
			},
		},
	)

	require.NoError(t, err)

	// Update the query schemas resource
	builder.UpdateSchemaDefinition(t, "../../src/static/schema/dataquery.json")
}
