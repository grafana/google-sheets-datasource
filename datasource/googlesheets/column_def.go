package googlesheets

import (
	"fmt"
	"strings"

	"github.com/araddon/dateparse"
	"google.golang.org/api/sheets/v4"
)

type columnDef struct {
	Name string
}

func getColumnType(columnIndex int, rows []*sheets.RowData) string {
	columnTypes := map[string]bool{}
	for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
		columnType := getType(rows[rowIndex].Values[columnIndex])
		columnTypes[columnType] = true
	}

	if len(columnTypes) == 1 {
		for key := range columnTypes {
			return key
		}
	}

	return "STRING"
}

func getColumnUnit(columnIndex int, rows []*sheets.RowData) string {
	columnUnits := map[string]bool{}
	for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
		columnUnit := getUnit(rows[rowIndex].Values[columnIndex])
		columnUnits[columnUnit] = true
	}

	if len(columnUnits) == 1 {
		for key := range columnUnits {
			return key
		}
	}

	return ""
}

func getType(cellData *sheets.CellData) string {
	if cellData.UserEnteredFormat.NumberFormat != nil {
		switch cellData.UserEnteredFormat.NumberFormat.Type {
		case "DATE", "DATE_TIME":
			return "TIME"
		case "NUMBER", "PERCENT", "CURRENCY":
			return "NUMBER"
		}
	}

	return "STRING"
}

func getUnit(cellData *sheets.CellData) string {
	if cellData.UserEnteredFormat.NumberFormat != nil {
		switch cellData.UserEnteredFormat.NumberFormat.Type {
		case "NUMBER":
			if strings.Contains(cellData.UserEnteredFormat.NumberFormat.Pattern, "$") {
				return "$"
			}
		case "PERCENT":
			return "%"
		case "CURRENCY":
			return "$"
		}
	}

	return ""
}

func getValue(cellData *sheets.CellData) (interface{}, error) {
	if cellData.UserEnteredFormat.NumberFormat != nil {
		switch cellData.UserEnteredFormat.NumberFormat.Type {
		case "DATE", "DATE_TIME":
			time, err := dateparse.ParseLocal(cellData.FormattedValue)
			if err != nil {
				return nil, fmt.Errorf("error while parsing date :", err.Error())
			}

			return &time, nil
		default:
			return &cellData.EffectiveValue.NumberValue, nil
		}
	}

	return &cellData.FormattedValue, nil
}
