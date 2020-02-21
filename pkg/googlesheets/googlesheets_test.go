package googlesheets

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/sheets/v4"
)

type fakeClient struct {
}

func (f *fakeClient) GetSpreadsheet(spreadSheetID string, sheetRange string, includeGridData bool) (*sheets.Spreadsheet, error) {
	sheet, _ := loadTestSheet("./testdata/mixed-data.json")
	return sheet, nil
}

func loadTestSheet(path string) (*sheets.Spreadsheet, error) {
	var data *sheets.Spreadsheet

	jsonBody, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBody, &data)
	sheet := data

	return sheet, nil
}

func TestGooglesheets(t *testing.T) {
	t.Run("getUniqueColumnName", func(t *testing.T) {
		t.Run("name is appended with number if not unique", func(t *testing.T) {
			columns := map[string]bool{"header": true, "name": true}
			name := getUniqueColumnName("header", 1, columns)
			assert.Equal(t, name, "header1")
		})

		t.Run("name becomes field + column index if header row is empty", func(t *testing.T) {
			columns := map[string]bool{}
			name := getUniqueColumnName("", 3, columns)
			assert.Equal(t, name, "Field 4")
		})
	})

	t.Run("getSpreadSheet", func(t *testing.T) {
		client := &fakeClient{}
		t.Run("spreadsheet is being cached", func(t *testing.T) {
			gsd := &GoogleSheets{
				Cache: cache.New(300*time.Second, 50*time.Second),
			}
			qm := QueryModel{Range: "A1:O", Spreadsheet: Spreadsheet{ID: "someid"}, CacheDurationSeconds: 10}

			assert.Equal(t, 0, gsd.Cache.ItemCount())

			_, meta, _ := gsd.getSpreadSheet(client, &qm)
			assert.False(t, meta["hit"].(bool))
			assert.Equal(t, 1, gsd.Cache.ItemCount())
		})

		t.Run("spreadsheet is not being cached if CacheDurationSeconds is 0", func(t *testing.T) {
			gsd := &GoogleSheets{
				Cache: cache.New(300*time.Second, 50*time.Second),
			}
			qm := QueryModel{Range: "A1:O", Spreadsheet: Spreadsheet{ID: "someid"}, CacheDurationSeconds: 0}

			assert.Equal(t, 0, gsd.Cache.ItemCount())

			_, meta, _ := gsd.getSpreadSheet(client, &qm)
			assert.False(t, meta["hit"].(bool))
			assert.Equal(t, 0, gsd.Cache.ItemCount())
		})
	})

	t.Run("transformSheetToDataFrame", func(t *testing.T) {
		sheet, _ := loadTestSheet("./testdata/mixed-data.json")
		gsd := &GoogleSheets{
			Cache: cache.New(300*time.Second, 50*time.Second),
			Logger: hclog.New(&hclog.LoggerOptions{
				Name:  "",
				Level: hclog.LevelFromString("DEBUG"),
			}),
		}
		qm := QueryModel{Range: "A1:O", Spreadsheet: Spreadsheet{ID: "someid"}, CacheDurationSeconds: 10}

		meta := make(map[string]interface{})
		frame, err := gsd.transformSheetToDataFrame(sheet.Sheets[0].Data[0], meta, "ref1", &qm)
		assert.Nil(t, err)
		assert.Equal(t, frame.Name, "ref1")

		t.Run("no of columns match", func(t *testing.T) {
			assert.Equal(t, len(frame.Fields), 16)
		})

		t.Run("no of rows match field length", func(t *testing.T) {
			for _, field := range frame.Fields {
				assert.Equal(t, field.Len(), len(sheet.Sheets[0].Data[0].RowData)-1)
			}
		})

		t.Run("meta is populated correctly", func(t *testing.T) {
			assert.Equal(t, meta["spreadsheetId"], qm.Spreadsheet.ID)
			assert.Equal(t, meta["range"], qm.Range)
		})

		t.Run("meta warnings is populated correctly", func(t *testing.T) {
			warnings := meta["warnings"].([]string)
			assert.Equal(t, 4, len(warnings))
			assert.Equal(t, "Multipe data types found in column MixedDataTypes. Using string data type", warnings[0])
			assert.Equal(t, "Multipe data types found in column MixedUnits. Using string data type", warnings[1])
			assert.Equal(t, "Multipe units found in column MixedUnits. Formatted value will be used", warnings[2])
			assert.Equal(t, "Multipe units found in column Mixed currencies. Formatted value will be used", warnings[3])
		})
	})
}
