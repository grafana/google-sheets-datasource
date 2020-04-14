package googlesheets

import (
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func findTimeField(frame *data.Frame) int {
	timeIndices := frame.TypeIndices(data.FieldTypeTime, data.FieldTypeNullableTime)
	if len(timeIndices) == 0 {
		return -1
	}
	return timeIndices[0]
}

func getExcelColumnName(columnNumber int) string {
	dividend := columnNumber
	columnName := ""
	var modulo int

	for dividend > 0 {
		modulo = ((dividend - 1) % 26)
		columnName = string(65+modulo) + columnName
		dividend = ((dividend - modulo) / 26)
	}

	return columnName
}
