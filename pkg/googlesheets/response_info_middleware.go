package googlesheets

import (
	"context"
	"net/http"

	"github.com/grafana/google-sheets-datasource/pkg/models"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
)

const ResponseInfoMiddlewareName = "response-info"

func ResponseInfoMiddleware() httpclient.Middleware {
	return httpclient.NamedMiddlewareFunc(ResponseInfoMiddlewareName, RoundTripper)
}

func RoundTripper(_ httpclient.Options, next http.RoundTripper) http.RoundTripper {
	return httpclient.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		res, err := next.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		res.Body = httpclient.CountBytesReader(res.Body, func(size int64) {
			backend.Logger.FromContext(req.Context()).Debug("Downstream response info", "bytes", size, "url", req.URL.String(), "retrieved", true)
		})
		return res, err
	})
}

func logIfNotAbleToRetrieveResponseInfo(ctx context.Context, settings models.DatasourceSettings) {
	if settings.AuthenticationType == authenticationTypeAPIKey && len(settings.APIKey) > 0 {
		backend.Logger.FromContext(ctx).Debug("Downstream response info", "retrieved", false)
	}
}
