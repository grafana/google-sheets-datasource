package ext

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/grafana/kindsys"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/klog/v2"
)

type ctxResourceKey struct{}
type ctxSubresourceKey struct{}

// WithResource adds a resource to the context.
func WithResource(ctx context.Context, r kindsys.Resource) context.Context {
	return context.WithValue(ctx, ctxResourceKey{}, r)
}

// WithSubresource adds a resource to the context.
func WithSubresource(ctx context.Context, subresource string) context.Context {
	return context.WithValue(ctx, ctxSubresourceKey{}, subresource)
}

// ResourceFromContext gets the resource from context
func ResourceFromContext(ctx context.Context) (kindsys.Resource, error) {
	u, ok := ctx.Value(ctxResourceKey{}).(kindsys.Resource)
	if ok && u != nil {
		return u, nil
	}
	return nil, fmt.Errorf("a Resource was not found in the context")
}

// SubresourceFromContext gets the resource from context
func SubresourceFromContext(ctx context.Context) (*string, error) {
	s, ok := ctx.Value(ctxSubresourceKey{}).(string)
	if ok {
		return &s, nil
	}
	return nil, fmt.Errorf("a Subresource was not found in the context")
}

// generic... would be in SDK
func SubresourceHandlerWrapper(upstream http.HandlerFunc, getter ResourceGetter) http.HandlerFunc {
	if upstream == nil {
		klog.Error("Nil handler passed to wrap with SubresourceHandlerWrapper")
		return http.NotFound
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		info, ok := request.RequestInfoFrom(ctx)
		if !ok {
			responsewriters.ErrorNegotiated(
				apierrors.NewInternalError(fmt.Errorf("no RequestInfo found in the context")),
				Codecs, schema.GroupVersion{}, writer, req,
			)
			return
		}

		r, err := getter(ctx, info.Namespace, info.Name)
		if err != nil {
			responsewriters.ErrorNegotiated(
				apierrors.NewInternalError(fmt.Errorf("could not get resource")),
				Codecs, schema.GroupVersion{}, writer, req,
			)
			return
		}

		ctx = WithResource(ctx, r)
		if info.Subresource != "" {
			subresource := info.Subresource
			// If we are using more parts than what Kubernetes expects for a subresource (say, it's a/b, not just a)
			// Parts for a resource request are: resource kind, resource name, subresource name
			// however, for subresource names that have their own segments, len(parts) > 3
			if len(info.Parts) > 3 {
				subresource = subresource + "/" + strings.Join(info.Parts[3:], "/")
			}
			ctx = WithSubresource(ctx, subresource)
		}

		upstream(writer, req.WithContext(ctx))
	}
}
