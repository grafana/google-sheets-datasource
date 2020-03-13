package utils

import (
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// FindTimeField gets the first time field index or -1
func FindTimeField(frame *data.Frame) int {
	for fieldIdx, f := range frame.Fields {
		ftype := f.Type()
		if ftype == data.FieldTypeTime || ftype == data.FieldTypeNullableTime {
			return fieldIdx
		}
	}
	return -1
}

// GetExcelColumnName -- get a name based on the index A,B,C...,AA,AB,... etc
func GetExcelColumnName(columnNumber int) string {
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
