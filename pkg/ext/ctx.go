package ext

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grafana/kindsys"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/klog/v2"
)

type ctxResourceKey struct{}

// WithResource adds a resource to the context.
func WithResource(ctx context.Context, r kindsys.Resource) context.Context {
	return context.WithValue(ctx, ctxResourceKey{}, r)
}

// ResourceFromContext gets the resource from context
func ResourceFromContext(ctx context.Context) (kindsys.Resource, error) {
	u, ok := ctx.Value(ctxResourceKey{}).(kindsys.Resource)
	if ok && u != nil {
		return u, nil
	}
	return nil, fmt.Errorf("a Resource was not found in the context")
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

		// Attach the resource to the context
		upstream(writer, req.WithContext(WithResource(ctx, r)))
	}
}
