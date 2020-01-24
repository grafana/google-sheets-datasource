package scenario

import (
	"net/http"
	"time"
)

var httpClient = http.Client{
	Timeout:   time.Second * 30,
	Transport: newAddUserAgent(nil),
}

type addUserAgent struct {
	T http.RoundTripper
}

func (aua *addUserAgent) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", "Grafana-Test-Datasource/Alpha")
	return aua.T.RoundTrip(req)
}

func newAddUserAgent(T http.RoundTripper) *addUserAgent {
	if T == nil {
		T = http.DefaultTransport
	}
	return &addUserAgent{T}
}
