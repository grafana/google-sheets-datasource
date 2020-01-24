package scenario

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	ds "github.com/grafana/grafana-plugin-sdk-go/backend"
	df "github.com/grafana/grafana-plugin-sdk-go/dataframe"
)

const typeCSVWaveQuery queryType = "csvWave"

type csvWaveQuery struct {
	baseQuery
	TimeStep  int64          `json:"timeStep"`
	CSVValues csvInputValues `json:"csvValues"`
	Shift     int64          `json:"shift"`
	Phase     int64          `json:"phase"`
}

type csvInputValues []*float64

func (cw *csvInputValues) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.TrimRight(strings.TrimSpace(s), ",") // Strip Trailing Comma
	sVals := strings.Split(s, ",")
	floats := make([]*float64, len(sVals))
	for i, rawValue := range sVals {
		var fv *float64
		rawValue = strings.TrimSpace(rawValue)
		switch strings.ToLower(rawValue) {
		case "nan":
			nan := math.NaN()
			fv = &nan
		case "null":
			break
		default:
			f, err := strconv.ParseFloat(rawValue, 64)
			if err != nil {
				return err
			}
			fv = &f
		}
		floats[i] = fv
	}
	*cw = floats
	return nil
}

func (cw csvWaveQuery) Execute(_ context.Context) ([]*df.Frame, error) {
	if cw.Shift >= cw.TimeStep {
		return nil, fmt.Errorf("Shift must be smaller than Time Step")
	}
	valuesLen := int64(len(cw.CSVValues))
	getValue := func(mod int64) (*float64, error) {
		var i int64
		for i = 0; i < valuesLen; i++ {
			if mod-i*cw.TimeStep == cw.Shift {
				phaseIndex := i
				if cw.Phase > 0 {
					phaseIndex = (i + cw.Phase) % valuesLen
				}
				return cw.CSVValues[phaseIndex], nil
			}
		}
		return nil, fmt.Errorf("csvWave: should not be here")
	}
	frame, err := predictableSeries(cw.TimeRange, cw.TimeStep, valuesLen, cw.Shift, cw.MaxDataPoints, getValue)
	if err != nil {
		return nil, err
	}
	return []*df.Frame{
		frame,
	}, nil
}

func predictableSeries(timeRange ds.TimeRange, timeStep, length, shift, maxPoints int64, getValue func(mod int64) (*float64, error)) (*df.Frame, error) {
	frame := df.New("csv_wave",
		df.NewField("time", nil, []time.Time{}),
		df.NewField("value", nil, []*float64{}),
	)
	addPoint := func(t time.Time, f *float64) {
		frame.Fields[0].Vector.Append(t)
		frame.Fields[1].Vector.Append(f)
	}

	from := timeRange.From.Unix()
	to := timeRange.To.Unix()

	timeCursor := from - (from % timeStep)
	timeCursor += shift

	wavePeriod := timeStep * length

	for i := int64(0); i < maxPoints && timeCursor < to; i++ {
		f, err := getValue(timeCursor % wavePeriod)
		if err != nil {
			return nil, err
		}
		t := time.Unix(timeCursor, 0)
		addPoint(t, f)
		timeCursor += timeStep
	}

	return frame, nil
}
