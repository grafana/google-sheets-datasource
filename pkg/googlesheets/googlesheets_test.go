package googlesheets

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/grafana/google-sheets-datasource/pkg/models"

	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
)

type fakeClient struct {
	mock.Mock
}

func (f *fakeClient) GetSpreadsheet(ctx context.Context, spreadSheetID string, sheetRange string, includeGridData bool) (*sheets.Spreadsheet, error) {
	args := f.Called(ctx, spreadSheetID, sheetRange, includeGridData)
	if spreadsheet, ok := args.Get(0).(*sheets.Spreadsheet); ok {
		return spreadsheet, args.Error(1)
	}
	return nil, args.Error(1)
}

func loadTestSheet(path string) (*sheets.Spreadsheet, error) {
	jsonBody, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var sheet sheets.Spreadsheet
	if err := json.Unmarshal(jsonBody, &sheet); err != nil {
		return nil, err
	}

	return &sheet, nil
}

func TestGooglesheets(t *testing.T) {
	t.Run("getUniqueColumnName", func(t *testing.T) {
		t.Run("name is appended with number if not unique", func(t *testing.T) {
			columns := map[string]bool{"header": true, "name": true}
			name := getUniqueColumnName("header", 1, columns)
			assert.Equal(t, "header1", name)
		})

		t.Run("name becomes Field + column index if header row is empty", func(t *testing.T) {
			columns := map[string]bool{}
			name := getUniqueColumnName("", 3, columns)
			assert.Equal(t, "Field 4", name)
		})
	})

	t.Run("getSheetData", func(t *testing.T) {
		t.Run("spreadsheets get cached", func(t *testing.T) {
			client := &fakeClient{}
			qm := models.QueryModel{Range: "A1:O", Spreadsheet: "someId", CacheDurationSeconds: 10}
			gsd := &GoogleSheets{
				Cache: cache.New(300*time.Second, 50*time.Second),
			}
			require.Equal(t, 0, gsd.Cache.ItemCount())

			client.On("GetSpreadsheet", context.Background(), qm.Spreadsheet, qm.Range, true).Return(loadTestSheet("./testdata/mixed-data.json"))

			_, meta, err := gsd.getSheetData(context.Background(), client, &qm)
			require.NoError(t, err)

			assert.False(t, meta["hit"].(bool))
			assert.Equal(t, 1, gsd.Cache.ItemCount())

			_, meta, err = gsd.getSheetData(context.Background(), client, &qm)
			require.NoError(t, err)
			assert.True(t, meta["hit"].(bool))
			assert.Equal(t, 1, gsd.Cache.ItemCount())
			client.AssertExpectations(t)
		})

		t.Run("spreadsheets don't get cached if CacheDurationSeconds is 0", func(t *testing.T) {
			client := &fakeClient{}
			qm := models.QueryModel{Range: "A1:O", Spreadsheet: "someId", CacheDurationSeconds: 0}
			gsd := &GoogleSheets{
				Cache: cache.New(300*time.Second, 50*time.Second),
			}
			require.Equal(t, 0, gsd.Cache.ItemCount())

			client.On("GetSpreadsheet", context.Background(), qm.Spreadsheet, qm.Range, true).Return(loadTestSheet("./testdata/mixed-data.json"))

			_, meta, err := gsd.getSheetData(context.Background(), client, &qm)
			require.NoError(t, err)

			assert.False(t, meta["hit"].(bool))
			assert.Equal(t, 0, gsd.Cache.ItemCount())
			client.AssertExpectations(t)
		})

		t.Run("api error 404", func(t *testing.T) {
			client := &fakeClient{}
			qm := &models.QueryModel{
				Spreadsheet:          "spreadsheet-id",
				Range:                "Sheet1!A1:B2",
				CacheDurationSeconds: 60,
			}
			gsd := &GoogleSheets{
				Cache: cache.New(300*time.Second, 50*time.Second),
			}
			client.On("GetSpreadsheet", context.Background(), qm.Spreadsheet, qm.Range, true).Return(&sheets.Spreadsheet{}, &googleapi.Error{
				Code:    404,
				Message: "Not found",
			})

			_, _, err := gsd.getSheetData(context.Background(), client, qm)

			assert.Error(t, err)
			assert.Equal(t, "spreadsheet not found", err.Error())
			client.AssertExpectations(t)
		})

		t.Run("error other than 404", func(t *testing.T) {
			client := &fakeClient{}
			qm := &models.QueryModel{
				Spreadsheet:          "spreadsheet-id",
				Range:                "Sheet1!A1:B2",
				CacheDurationSeconds: 60,
			}
			gsd := &GoogleSheets{
				Cache: cache.New(300*time.Second, 50*time.Second),
			}
			client.On("GetSpreadsheet", context.Background(), qm.Spreadsheet, qm.Range, true).Return(&sheets.Spreadsheet{}, &googleapi.Error{
				Code:    403,
				Message: "Forbidden",
			})

			_, _, err := gsd.getSheetData(context.Background(), client, qm)

			assert.Error(t, err)
			assert.Equal(t, "google API Error 403", err.Error())

			client.AssertExpectations(t)
		})

		t.Run("context canceled", func(t *testing.T) {
			client := &fakeClient{}
			qm := &models.QueryModel{
				Spreadsheet:          "spreadsheet-id",
				Range:                "Sheet1!A1:B2",
				CacheDurationSeconds: 60,
			}
			gsd := &GoogleSheets{
				Cache: cache.New(300*time.Second, 50*time.Second),
			}
			client.On("GetSpreadsheet", context.Background(), qm.Spreadsheet, qm.Range, true).Return(&sheets.Spreadsheet{}, context.Canceled)

			_, _, err := gsd.getSheetData(context.Background(), client, qm)

			assert.Error(t, err)
			assert.Equal(t, context.Canceled.Error(), err.Error())
			assert.True(t, backend.IsDownstreamError(err))

			client.AssertExpectations(t)
		})

		t.Run("timeout", func(t *testing.T) {
			client := &fakeClient{}
			qm := &models.QueryModel{
				Spreadsheet:          "spreadsheet-id",
				Range:                "Sheet1!A1:B2",
				CacheDurationSeconds: 60,
			}
			gsd := &GoogleSheets{
				Cache: cache.New(300*time.Second, 50*time.Second),
			}

			client.On("GetSpreadsheet", context.Background(), qm.Spreadsheet, qm.Range, true).Return(&sheets.Spreadsheet{}, &net.OpError{Err: context.DeadlineExceeded})

			_, _, err := gsd.getSheetData(context.Background(), client, qm)

			assert.Error(t, err)
			assert.True(t, backend.IsDownstreamError(err))

			client.AssertExpectations(t)
		})

		t.Run("oauth invalid grant", func(t *testing.T) {
			client := &fakeClient{}
			qm := &models.QueryModel{
				Spreadsheet:          "spreadsheet-id",
				Range:                "Sheet1!A1:B2",
				CacheDurationSeconds: 60,
			}
			gsd := &GoogleSheets{
				Cache: cache.New(300*time.Second, 50*time.Second),
			}

			// Simulated oauth2.RetrieveError
			retrieveErr := &oauth2.RetrieveError{
				Response: &http.Response{
					Status:     "400 Bad Request",
					StatusCode: 400,
				},
				Body: []byte(`{"error":"invalid_grant","error_description":"Invalid grant: account not found"}`),
			}

			// Simulated *url.Error wrapping the retrieveErr
			urlErr := &url.Error{
				Op:  "Get",
				URL: "https://sheets.googleapis.com/v4/spreadsheets/...",
				Err: retrieveErr,
			}

			client.On("GetSpreadsheet", context.Background(), qm.Spreadsheet, qm.Range, true).Return(&sheets.Spreadsheet{}, urlErr)

			_, _, err := gsd.getSheetData(context.Background(), client, qm)

			assert.Error(t, err)
			assert.True(t, backend.IsDownstreamError(err))

			client.AssertExpectations(t)
		})

		t.Run("error that doesn't have message property", func(t *testing.T) {
			client := &fakeClient{}
			qm := &models.QueryModel{
				Spreadsheet:          "spreadsheet-id",
				Range:                "Sheet1!A1:B2",
				CacheDurationSeconds: 60,
			}
			gsd := &GoogleSheets{
				Cache: cache.New(300*time.Second, 50*time.Second),
			}

			client.On("GetSpreadsheet", context.Background(), qm.Spreadsheet, qm.Range, true).Return(&sheets.Spreadsheet{}, &googleapi.Error{
				Message: "",
			})

			_, _, err := gsd.getSheetData(context.Background(), client, qm)

			assert.Error(t, err)
			assert.Equal(t, "unknown API error", err.Error())

			client.AssertExpectations(t)
		})
	})

	t.Run("transformSheetToDataFrame", func(t *testing.T) {
		sheet, err := loadTestSheet("./testdata/mixed-data.json")
		require.NoError(t, err)

		gsd := &GoogleSheets{
			Cache: cache.New(300*time.Second, 50*time.Second),
		}
		qm := models.QueryModel{Range: "A1:O", Spreadsheet: "someId", CacheDurationSeconds: 10}

		meta := make(map[string]any)
		frame, err := gsd.transformSheetToDataFrame(context.Background(), sheet.Sheets[0].Data[0], meta, "ref1", &qm)
		require.NoError(t, err)
		require.Equal(t, "ref1", frame.Name)

		t.Run("no of columns match", func(t *testing.T) {
			assert.Equal(t, 16, len(frame.Fields))
		})

		t.Run("no of rows matches field length", func(t *testing.T) {
			for _, field := range frame.Fields {
				assert.Equal(t, len(sheet.Sheets[0].Data[0].RowData)-1, field.Len())
			}
		})

		t.Run("meta is populated correctly", func(t *testing.T) {
			assert.Equal(t, qm.Spreadsheet, meta["spreadsheetId"])
			assert.Equal(t, qm.Range, meta["range"])
		})

		t.Run("meta warnings field is populated correctly", func(t *testing.T) {
			warnings, ok := meta["warnings"].([]string)
			require.True(t, ok)
			assert.Equal(t, 3, len(warnings))
			assert.Equal(t, "Multiple data types found in column \"MixedDataTypes\". Using string data type", warnings[0])
			assert.Equal(t, "Multiple units found in column \"MixedUnits\". Formatted value will be used", warnings[1])
			assert.Equal(t, "Multiple units found in column \"Mixed currencies\". Formatted value will be used", warnings[2])
			// assert.Equal(t, "Multiple data types found in column \"MixedUnits\". Using string data type", warnings[2])
		})
	})

	t.Run("query single cell", func(t *testing.T) {
		sheet, err := loadTestSheet("./testdata/single-cell.json")
		require.NoError(t, err)

		gsd := &GoogleSheets{
			Cache: cache.New(300*time.Second, 50*time.Second),
		}
		qm := models.QueryModel{Range: "A2", Spreadsheet: "someId", CacheDurationSeconds: 10}

		meta := make(map[string]any)
		frame, err := gsd.transformSheetToDataFrame(context.Background(), sheet.Sheets[0].Data[0], meta, "ref1", &qm)
		require.NoError(t, err)
		require.Equal(t, "ref1", frame.Name)

		t.Run("single field", func(t *testing.T) {
			assert.Equal(t, 1, len(frame.Fields))
		})

		t.Run("single row", func(t *testing.T) {
			for _, field := range frame.Fields {
				assert.Equal(t, 1, field.Len())
			}
		})

		t.Run("single value", func(t *testing.T) {
			strVal, ok := frame.Fields[0].At(0).(*string)
			require.True(t, ok)
			require.NotNil(t, strVal)
			assert.Equal(t, "🌭", *strVal)
		})
	})

	t.Run("column id formatting", func(t *testing.T) {
		require.Equal(t, "A", getExcelColumnName(1))
		require.Equal(t, "B", getExcelColumnName(2))
		require.Equal(t, "AH", getExcelColumnName(34))
		require.Equal(t, "BN", getExcelColumnName(66))
		require.Equal(t, "ZW", getExcelColumnName(699))
		// cspell:disable-next-line
		require.Equal(t, "AJIL", getExcelColumnName(24582))
	})
}
