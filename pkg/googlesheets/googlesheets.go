package googlesheets

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/grafana/google-sheets-datasource/pkg/models"

	"github.com/araddon/dateparse"
	"github.com/davecgh/go-spew/spew"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/errorsource"
	"github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
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
		dr = errorsource.Response(err)
		dr.Error = fmt.Errorf("unable to create Google API client: %w", err)
		return
	}

	// This result may be cached
	sheetData, meta, err := gs.getSheetData(ctx, client, qm)
	if err != nil {
		dr = errorsource.Response(err)
		return
	}

	frame, err := gs.transformSheetToDataFrame(ctx, sheetData, meta, refID, qm)
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

	fileNames := make(map[string]string, len(files))
	for _, i := range files {
		fileNames[i.Id] = i.Name
	}

	return fileNames, nil
}

// getSheetData gets grid data corresponding to a spreadsheet.
func (gs *GoogleSheets) getSheetData(ctx context.Context, client client, qm *models.QueryModel) (*sheets.GridData, map[string]any, error) {
	logger := backend.Logger.FromContext(ctx)
	cacheKey := qm.Spreadsheet + qm.Range
	if item, expires, found := gs.Cache.GetWithExpiration(cacheKey); found && qm.CacheDurationSeconds > 0 {
		if gridData, ok := item.(*sheets.GridData); ok {
			return gridData, map[string]any{
				"hit":     true,
				"expires": expires.Unix(),
			}, nil
		}
		return nil, nil, errors.New("invalid cache item not type of *sheets.GridData")
	}
	result, err := client.GetSpreadsheet(ctx, qm.Spreadsheet, qm.Range, true)
	if err != nil {
		if apiErr, ok := err.(*googleapi.Error); ok {
			// Handle API-specific errors
			// We use ErrorSourceFromHTTPStatus to determine error source based on HTTP status code
			if apiErr.Code == 404 {
				errWithSource := errorsource.DownstreamError(errors.New("spreadsheet not found"), false)
				return nil, nil, errWithSource
			}
			if apiErr.Message != "" {
				logger.Warn("Google API Error: " + apiErr.Message)
				errWithSource := errorsource.SourceError(backend.ErrorSourceFromHTTPStatus(apiErr.Code), fmt.Errorf("google API Error %d", apiErr.Code), false)
				return nil, nil, errWithSource
			}
			errWithSource := errorsource.SourceError(backend.ErrorSourceFromHTTPStatus(apiErr.Code), errors.New("unknown API error"), false)
			logger.Warn(apiErr.Error())
			return nil, nil, errWithSource
		}

		if backend.IsDownstreamHTTPError(err) {
			errWithSource := errorsource.DownstreamError(err, false)
			return nil, nil, errWithSource
		}

		netErr, neErrOk := err.(net.Error)
		if neErrOk {
			var retrieveErr *oauth2.RetrieveError
			if errors.As(netErr, &retrieveErr) {
				errWithSource := errorsource.SourceError(backend.ErrorSourceFromHTTPStatus(retrieveErr.Response.StatusCode), err, false)
				return nil, nil, errWithSource
			}
		}

		logger.Warn("unknown error", "err", err)
		// This is an unknown error from the client - it might have error source middleware.
		// If not, it will be handled by the default error source - plugin error.
		return nil, nil, err
	}

	if result.Properties.TimeZone != "" {
		loc, err := time.LoadLocation(result.Properties.TimeZone)
		if err != nil {
			logger.Warn("could not load timezone from spreadsheet: %w", err)
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

func (gs *GoogleSheets) transformSheetToDataFrame(ctx context.Context, sheet *sheets.GridData, meta map[string]any, refID string, qm *models.QueryModel) (*data.Frame, error) {
	logger := backend.Logger.FromContext(ctx)
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
			logger.Debug(warning)
		}

		if column.HasMixedUnits() {
			warning := fmt.Sprintf("Multiple units found in column %q. Formatted value will be used", column.Header)
			warnings = append(warnings, warning)
			logger.Debug(warning)
		}
	}

	// We want to show the warnings only once per column
	warningsIncludeConverterErrorForColumns := make(map[int]bool, len(columns))
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
			if err != nil && !warningsIncludeConverterErrorForColumns[columnIndex] {
				logger.Debug("unsuccessful converting of cell data", "err", err)
				warnings = append(warnings, err.Error())
				warningsIncludeConverterErrorForColumns[columnIndex] = true
			}
		}
	}

	meta["warnings"] = warnings
	meta["spreadsheetId"] = qm.Spreadsheet
	meta["range"] = qm.Range
	frame.Meta = &data.FrameMeta{Custom: meta}
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
	if len(rows) < 1 {
		return []*ColumnDefinition{}, 0
	}
	headerRow := rows[0].Values

	columns := make([]*ColumnDefinition, 0, len(headerRow))
	columnMap := make(map[string]bool, len(headerRow))
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
