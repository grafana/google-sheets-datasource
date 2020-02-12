package googlesheets

import (
	"fmt"
	"sort"
	"strings"

	"github.com/araddon/dateparse"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/api/sheets/v4"
)

type columnDefinition struct {
	Header      string
	ColumnIndex int
	Type        string
	Unit        string
	Warning     string
}

func getColumnDefintions(rows []*sheets.RowData, logger hclog.Logger) []*columnDefinition {
	columnTypes := map[int]map[string]*columnDefinition{}
	for columnIndex := range rows[0].Values {
		columnTypes[columnIndex] = map[string]*columnDefinition{}
	}

	for rowIndex := 0; rowIndex < len(rows); rowIndex++ {
		for columnIndex, columnCell := range rows[rowIndex].Values {
			columnType := getType(columnCell)
			unit := getUnit(columnCell)
			columnTypes[columnIndex][columnType] = &columnDefinition{Type: columnType, Unit: unit}
		}
	}

	columns := []*columnDefinition{}
	for columnIndex, columnTypeMap := range columnTypes {
		columnName := rows[0].Values[columnIndex].FormattedValue
		var column *columnDefinition
		if len(columnTypeMap) == 1 {
			for _, c := range columnTypeMap {
				column = c
			}
		} else {
			//The column has different data types - fallback to string
			column = &columnDefinition{Type: "STRING", Warning: fmt.Sprint("Multipe data types found in column index %v. Using string data type", columnName)}
		}
		column.ColumnIndex = columnIndex
		column.Header = columnName
		columns = append(columns, column)
	}

	sort.Slice(columns, func(i, j int) bool {
		return columns[i].ColumnIndex < columns[j].ColumnIndex
	})

	return columns
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
