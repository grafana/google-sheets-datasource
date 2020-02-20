package googlesheets

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/sheets/v4"
)

func loadTestSheet(path string) (*sheets.GridData, error) {
	var data *sheets.Spreadsheet

	jsonBody, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBody, &data)
	sheet := data.Sheets[0].Data[0]

	return sheet, nil
}

func TestGooglesheets(t *testing.T) {
	// sheet, _ := loadTestSheet("../testdata/mixed-data.json")

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
}
