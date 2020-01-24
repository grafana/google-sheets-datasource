package scenario

import (
	"context"
	"time"

	df "github.com/grafana/grafana-plugin-sdk-go/dataframe"
)

const typeWalkQuery queryType = "walk"

type walkQuery struct {
	SomeFloat float64 `json:"someFloat"`
	baseQuery
}

func (wq walkQuery) Execute(ctx context.Context) ([]*df.Frame, error) {
	frame := df.New("walk",
		df.NewField("Time", nil, []time.Time{wq.TimeRange.From}),
		df.NewField("Value", nil, []*float64{&wq.SomeFloat}))
	frame.RefID = wq.RefID
	return []*df.Frame{frame}, nil
}
