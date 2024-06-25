package googlesheets

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/google-sheets-datasource/pkg/models"

	"github.com/araddon/dateparse"
	"github.com/davecgh/go-spew/spew"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/patrickmn/go-cache"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
)

// GoogleSheets provides an interface to the Google Sheets API.
type GoogleSheets struct {
	Cache *cache.Cache
}

// Query queries a spreadsheet and returns a corresponding data frame.
func (gs *GoogleSheets) Query(ctx context.Context, refID string, qm *models.QueryModel, config models.DatasourceSettings, timeRange backend.TimeRange) (dr backend.DataResponse) {
	client, err := NewGoogleClient(ctx, config)
	if err != nil {
		dr.Error = fmt.Errorf("unable to create Google API client: %w", err)
		return
	}

	// This result may be cached
	sheetData, meta, err := gs.getSheetData(client, qm)
	if err != nil {
		dr.Error = err
		return
	}

	frame, err := gs.transformSheetToDataFrame(sheetData, meta, refID, qm)
	if err != nil {
		dr.Error = err
		return
	}
	if frame == nil {
		return
	}
	if qm.UseTimeFilter {
		timeIndex := findTimeField(frame)
		if timeIndex >= 0 {
			frame, dr.Error = frame.FilterRowsByField(timeIndex, func(i any) (bool, error) {
				val, ok := i.(*time.Time)
				if !ok {
					return false, fmt.Errorf("invalid time column: %s", spew.Sdump(i))
				}
				if val == nil || val.Before(timeRange.From) || val.After(timeRange.To) {
					return false, nil
				}
				return true, nil
			})
		}
	}
	dr.Frames = append(dr.Frames, frame)
	return
}

// GetSpreadsheets gets spreadsheets from the Google API.
func (gs *GoogleSheets) GetSpreadsheets(ctx context.Context, config models.DatasourceSettings) (map[string]string, error) {
	client, err := NewGoogleClient(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google API client: %w", err)
	}

	files, err := client.GetSpreadsheetFiles()
	if err != nil {
		return nil, err
	}

	fileNames := map[string]string{}
	for _, i := range files {
		fileNames[i.Id] = i.Name
	}

	return fileNames, nil
}

// getSheetData gets grid data corresponding to a spreadsheet.
func (gs *GoogleSheets) getSheetData(client client, qm *models.QueryModel) (*sheets.GridData, map[string]any, error) {
	cacheKey := qm.Spreadsheet + qm.Range
	if item, expires, found := gs.Cache.GetWithExpiration(cacheKey); found && qm.CacheDurationSeconds > 0 {
		return item.(*sheets.GridData), map[string]any{
			"hit":     true,
			"expires": expires.Unix(),
		}, nil
	}

	result, err := client.GetSpreadsheet(qm.Spreadsheet, qm.Range, true)
	if err != nil {
		if apiErr, ok := err.(*googleapi.Error); ok {
			// Handle API-specific errors
			if apiErr.Code == 404 {
				return nil, nil, errors.New("spreadsheet not found")
			}
			if apiErr.Message != "" {
				log.DefaultLogger.Error("Google API Error: " + apiErr.Message)
				return nil, nil, fmt.Errorf("Google API Error %d", apiErr.Code)
			}
			log.DefaultLogger.Error(apiErr.Error())
			return nil, nil, errors.New("unknown API error")
		}
		log.DefaultLogger.Error("unknown error", "err", err)
		return nil, nil, errors.New("unknown error")
	}

	if result.Properties.TimeZone != "" {
		loc, err := time.LoadLocation(result.Properties.TimeZone)
		if err != nil {
			log.DefaultLogger.Warn("could not load timezone from spreadsheet: %w", err)
		} else {
			time.Local = loc
		}
	}

	sheetData := result.Sheets[0].Data[0]
	if qm.CacheDurationSeconds > 0 {
		gs.Cache.Set(cacheKey, sheetData, time.Duration(qm.CacheDurationSeconds)*time.Second)
	}

	return sheetData, map[string]any{"hit": false}, nil
}

func (gs *GoogleSheets) transformSheetToDataFrame(sheet *sheets.GridData, meta map[string]any, refID string, qm *models.QueryModel) (*data.Frame, error) {
	columns, start := getColumnDefinitions(sheet.RowData)
	warnings := []string{}

	converters := make([]data.FieldConverter, len(columns))
	for i, column := range columns {
		fc, ok := converterMap[column.GetType()]
		if !ok {
			return nil, fmt.Errorf("unknown column type: %s", column.GetType())
		}
		converters[i] = fc
	}

	inputConverter, err := data.NewFrameInputConverter(converters, len(sheet.RowData)-start)
	if err != nil {
		return nil, err
	}
	frame := inputConverter.Frame
	frame.RefID = refID
	frame.Name = refID // TODO: should set the name from metadata

	for i, column := range columns {
		field := frame.Fields[i]
		field.Name = column.Header
		field.Config = &data.FieldConfig{
			DisplayName: column.Header,
			Unit:        column.GetUnit(),
		}
		if column.HasMixedTypes() {
			warning := fmt.Sprintf("Multiple data types found in column %q. Using string data type", column.Header)
			warnings = append(warnings, warning)
			log.DefaultLogger.Warn(warning)
		}

		if column.HasMixedUnits() {
			warning := fmt.Sprintf("Multiple units found in column %q. Formatted value will be used", column.Header)
			warnings = append(warnings, warning)
			log.DefaultLogger.Warn(warning)
		}
	}

	for rowIndex := start; rowIndex < len(sheet.RowData); rowIndex++ {
		for columnIndex, cellData := range sheet.RowData[rowIndex].Values {
			if columnIndex >= len(columns) {
				continue
			}

			// Skip any empty values
			if cellData.FormattedValue == "" {
				continue
			}

			err := inputConverter.Set(columnIndex, rowIndex-start, cellData)
			if err != nil {
				warnings = append(warnings, err.Error())
			}
		}
	}

	meta["warnings"] = warnings
	meta["spreadsheetId"] = qm.Spreadsheet
	meta["range"] = qm.Range
	frame.Meta = &data.FrameMeta{Custom: meta}
	log.DefaultLogger.Debug("frame.Meta: %s", spew.Sdump(frame.Meta))
	return frame, nil
}

// timeConverter handles sheets TIME column types.
var timeConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableTime,
	Converter: func(i any) (any, error) {
		var t *time.Time
		cellData, ok := i.(*sheets.CellData)
		if !ok {
			return t, fmt.Errorf("expected type *sheets.CellData, but got %T", i)
		}
		parsedTime, err := dateparse.ParseLocal(cellData.FormattedValue)
		if err != nil {
			return t, fmt.Errorf("error while parsing date '%v'", cellData.FormattedValue)
		}
		return &parsedTime, nil
	},
}

// stringConverter handles sheets STRING column types.
var stringConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableString,
	Converter: func(i any) (any, error) {
		var s *string
		cellData, ok := i.(*sheets.CellData)
		if !ok {
			return s, fmt.Errorf("expected type *sheets.CellData, but got %T", i)
		}
		return &cellData.FormattedValue, nil
	},
}

// numberConverter handles sheets STRING column types.
var numberConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableFloat64,
	Converter: func(i any) (any, error) {
		cellData, ok := i.(*sheets.CellData)
		if !ok {
			return nil, fmt.Errorf("expected type *sheets.CellData, but got %T", i)
		}
		return cellData.EffectiveValue.NumberValue, nil
	},
}

// converterMap is a map sheets.ColumnType to fieldConverter and
// is used to create a data.FrameInputConverter for a returned sheet.
var converterMap = map[ColumnType]data.FieldConverter{
	"TIME":   timeConverter,
	"STRING": stringConverter,
	"NUMBER": numberConverter,
}

func getUniqueColumnName(formattedName string, columnIndex int, columns map[string]bool) string {
	name := formattedName
	if name == "" {
		name = fmt.Sprintf("Field %d", columnIndex+1)
	}

	nameExist := true
	counter := 1
	for nameExist {
		if _, exist := columns[name]; exist {
			name = fmt.Sprintf("%s%d", formattedName, counter)
			counter++
		} else {
			nameExist = false
		}
	}

	return name
}

func getColumnDefinitions(rows []*sheets.RowData) ([]*ColumnDefinition, int) {
	columns := []*ColumnDefinition{}
	columnMap := map[string]bool{}
	if len(rows) < 1 {
		return columns, 0
	}
	headerRow := rows[0].Values

	start := 0
	if len(rows) > 1 {
		start = 1
		for columnIndex, headerCell := range headerRow {
			name := getUniqueColumnName(strings.TrimSpace(headerCell.FormattedValue), columnIndex, columnMap)
			columnMap[name] = true
			columns = append(columns, NewColumnDefinition(name, columnIndex))
		}
	} else {
		for columnIndex := range headerRow {
			name := getUniqueColumnName("", columnIndex, columnMap)
			columnMap[name] = true
			columns = append(columns, NewColumnDefinition(name, columnIndex))
		}
	}

	// Check the types for each column
	for rowIndex := start; rowIndex < len(rows); rowIndex++ {
		for _, column := range columns {
			if column.ColumnIndex < len(rows[rowIndex].Values) {
				column.CheckCell(rows[rowIndex].Values[column.ColumnIndex])
			}
		}
	}

	return columns, start
}
