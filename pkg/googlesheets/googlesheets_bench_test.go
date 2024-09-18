package googlesheets

import (
	"context"
	"testing"

	"github.com/grafana/google-sheets-datasource/pkg/models"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stretchr/testify/require"
)

// To avoid compiler optimizations eliminating the function under test
// we are storing the result to a package level variable
var Frame *data.Frame

func BenchmarkTransformMixedSheetToDataFrame(b *testing.B) {
	sheet, err := loadTestSheet("./testdata/mixed-data.json")
	require.NoError(b, err)
	gsd := &GoogleSheets{}
	qm := models.QueryModel{Spreadsheet: "someId"}
	meta := make(map[string]any)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		frame, err := gsd.transformSheetToDataFrame(context.Background(), sheet.Sheets[0].Data[0], meta, "ref1", &qm)
		require.NoError(b, err)
		Frame = frame
	}
}

func BenchmarkTransformMixedSheetWithInvalidDateToDataFrame(b *testing.B) {
	sheet, err := loadTestSheet("./testdata/invalid-date-time.json")
	require.NoError(b, err)
	gsd := &GoogleSheets{}
	qm := models.QueryModel{Spreadsheet: "someId"}
	meta := make(map[string]any)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		frame, err := gsd.transformSheetToDataFrame(context.Background(), sheet.Sheets[0].Data[0], meta, "ref1", &qm)
		require.NoError(b, err)
		Frame = frame
	}
}
