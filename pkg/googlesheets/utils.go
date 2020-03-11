package googlesheets

import (
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func findTimeField(frame *data.Frame) int {
	for fieldIdx, f := range frame.Fields {
		// The vector type is time.Time!
		if f.Name == "time" { // DOOH!!  how do we check if it is a time field?
			return fieldIdx
		}
	}
	return -1
}

func filterByTime(frame *data.Frame, timeRange backend.TimeRange) *data.Frame {
	timeIndex := findTimeField(frame)
	if timeIndex < 0 {
		return frame // no time field... just leave it alone
	}

	timeField := frame.Fields[timeIndex]
	length := timeField.Len()
	filteredIndex := make([]int, length)
	filteredLen := 0

	for i := 0; i < length; i++ {
		val := timeField.At(i).(time.Time)
		if val.Before(timeRange.From) || val.After(timeRange.To) {
			continue
		}
		filteredIndex[filteredLen] = i
		filteredLen++ // increment the length
	}

	if filteredLen != length {
		// TODO -- create a new frame keeping all fields, but shorter vectors
		fmt.Println("TODO... actually filter")
	}
	return frame
}
