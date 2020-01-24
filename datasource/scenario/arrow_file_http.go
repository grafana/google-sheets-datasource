package scenario

import (
	"context"
	"io/ioutil"
	"net/http"

	df "github.com/grafana/grafana-plugin-sdk-go/dataframe"
)

const typeArrowFileQuery queryType = "arrowFile"

type arrowFileQuery struct {
	URL string `json:"url"`
	baseQuery
}

func (af arrowFileQuery) Execute(ctx context.Context) ([]*df.Frame, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, af.URL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	frame, err := df.UnmarshalArrow(body)
	if err != nil {
		return nil, err
	}
	return []*df.Frame{
		frame,
	}, nil
}
