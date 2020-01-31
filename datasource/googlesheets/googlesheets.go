package googlesheets

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	df "github.com/grafana/grafana-plugin-sdk-go/dataframe"
	"github.com/hashicorp/go-hclog"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func createService(ctx context.Context, config *GoogleSheetConfig) (*sheets.Service, error) {
	if config.AuthType == "none" {
		return sheets.NewService(ctx, option.WithAPIKey(config.ApiKey))
	}

	jwtConfig, err := google.JWTConfigFromJSON([]byte(config.JwtFile), "https://www.googleapis.com/auth/spreadsheets.readonly", "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, fmt.Errorf("Error parsin JWT file: %v", err)
	}

	client := jwtConfig.Client(ctx)

	return sheets.New(client)
}

func getTypeDefaultValue(t string) (interface{}) {
	switch t {
	case "time":
		return nil
	case "float64":
		return 0.0
	default:
		return ""
	}
}

// Query function
func Query(ctx context.Context, refID string, qm *QueryModel, config *GoogleSheetConfig, logger hclog.Logger) (*df.Frame, error) {
	srv, err := createService(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Unable to create service: %v", err.Error())
	}

	result, err := srv.Spreadsheets.Get(qm.SpreadsheetID).Ranges(qm.Range).IncludeGridData(true).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to get spreadsheet: %v", err.Error())
	}

	sheet := result.Sheets[0].Data[0]

	// Ensure all cells in a column have same type
	columnDef := map[int]string{}
	for rIndex, row := range sheet.RowData {
		if rIndex == 0 {
			continue
		}

		for cIndex, sCell := range row.Values {
			cell := &cellParser{cell: sCell}
			if val, exists := columnDef[cIndex]; exists {
				t := cell.GetType()
				if t != "" && val != t {
					return nil, fmt.Errorf("Column %v contains different data types. Found in row %v", cIndex, rIndex)
				}
			} else {
				columnDef[cIndex] = cell.GetType()
			}
		}
	}

	// Create fields using header name
	fields := []*df.Field{}
	// for _, sCell := range sheet.RowData[0].Values {
	for columnIndex, t := range columnDef {
		// cell := &cellParser{cell: sCell}
		cellParser := &cellParser{cell: sheet.RowData[0].Values[columnIndex]}
		switch t {
		case "time":
			fields = append(fields, df.NewField(cellParser.cell.EffectiveValue.StringValue, nil, []time.Time{}))
		case "float64":
			fields = append(fields, df.NewField(cellParser.cell.EffectiveValue.StringValue, nil, []float64{}))
		default:
			fields = append(fields, df.NewField(cellParser.cell.EffectiveValue.StringValue, nil, []string{}))
		}
	}

	frame := df.New(qm.Range, fields...)
	frame.RefID = refID
	for rowIndex := 1; rowIndex < len(sheet.RowData); rowIndex++ {
		for columnIndex, t := range columnDef {
			if columnIndex+1 > len(sheet.RowData[rowIndex].Values) {
				frame.Fields[columnIndex].Vector.Append(getTypeDefaultValue(t))
			} else {
				cellParser := &cellParser{cell: sheet.RowData[rowIndex].Values[columnIndex]}
				cellType := cellParser.GetType()
				if cellType == "" {
					frame.Fields[columnIndex].Vector.Append(getTypeDefaultValue(t))
				} else {
				v, err := cellParser.getValue()
				if err != nil {
					return nil, fmt.Errorf("Could not parse value: ", err.Error())
				}
				logger.Debug("type: " + t)
				logger.Debug("value: " + sheet.RowData[rowIndex].Values[columnIndex].FormattedValue)
				logger.Debug("janneb: " + spew.Sdump(v))
				frame.Fields[columnIndex].Vector.Append(v)
				}
			}
		}
	}

	return frame, nil
}

// TestAPI function
func TestAPI(ctx context.Context, config *GoogleSheetConfig) (*df.Frame, error) {
	_, err := createService(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Invalid datasource configuration: %s", err)
	}

	return df.New("TestAPI"), nil
}
