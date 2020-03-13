package googlesheets

import (
	"google.golang.org/api/sheets/v4"
)

type client interface {
	GetSpreadsheet(spreadSheetID string, sheetRange string, includeGridData bool) (*sheets.Spreadsheet, error)
}
