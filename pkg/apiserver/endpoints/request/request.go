// SPDX-License-Identifier: AGPL-3.0-only

package request

import (
	"context"
	"net/http"

	"k8s.io/apiserver/pkg/endpoints/handlers"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
)

type outputMediaType int

const outputMediaKey outputMediaType = iota

func WithOutputMediaType(ctx context.Context, req *http.Request, scope *handlers.RequestScope) context.Context {
	outputMedia, _, err := negotiation.NegotiateOutputMediaType(req, scope.Serializer, scope)
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, outputMediaKey, outputMedia)
}

func OutputMediaTypeFrom(ctx context.Context) (negotiation.MediaTypeOptions, bool) {
	mt, ok := ctx.Value(outputMediaKey).(negotiation.MediaTypeOptions)
	return mt, ok
}
