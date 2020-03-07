package googlesheets

import (
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	df "github.com/grafana/grafana-plugin-sdk-go/dataframe"
)

func findTimeField(frame *df.Frame) int {
	for fieldIdx, f := range frame.Fields {
		// The vector type is time.Time!
		if f.Name == "time" { // DOOH!!  how do we check if it is a time field?
			return fieldIdx
		}
	}
	return -1
}

func filterByTime(frame *df.Frame, timeRange backend.TimeRange) *df.Frame {
	timeIndex := findTimeField(frame)
	if timeIndex < 0 {
		return frame // no time field... just leave it alone
	}

	timeVector := frame.Fields[timeIndex].Vector
	length := timeVector.Len()
	filteredIndex := make([]int, length)
	filteredLen := 0

	for i := 0; i < length; i++ {
		val := timeVector.At(i).(time.Time)
		if val.Before(timeRange.From) || val.After(timeRange.To) {
			continue
		}
		filteredIndex[filteredLen] = i
		filteredLen++ // increment the length
	}

	if filteredLen != length {
		// TODO -- create a new frame keeping all fields, but shorter vectors
		length = filteredLen // make linter happy
	}
	return frame
}
