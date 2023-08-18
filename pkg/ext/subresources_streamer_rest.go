package ext

import (
	"context"
	"fmt"

	v1 "github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets/v1"
	"github.com/grafana/google-sheets-datasource/pkg/apiserver/apihelpers"
	"github.com/grafana/kindsys"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	restclient "k8s.io/client-go/rest"
)

var _ rest.Storage = (*SubresourceStreamerREST)(nil)
var _ rest.Getter = (*SubresourceStreamerREST)(nil)
var _ rest.CreaterUpdater = (*SubresourceStreamerREST)(nil)

type SubresourceStreamerREST struct {
	RestConfig *restclient.Config
}

func (r *SubresourceStreamerREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	//TODO implement me
	panic("implement me")
}

func (r *SubresourceStreamerREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *SubresourceStreamerREST) New() runtime.Object {
	return &v1.Datasource{}
}

func (r *SubresourceStreamerREST) Destroy() {
}

func (r *SubresourceStreamerREST) Get(ctx context.Context, name string, _ *metav1.GetOptions) (runtime.Object, error) {
	streamer := &apihelpers.SubresourceStreamer{}

	info, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("could not get request info")
	}
	switch info.Subresource {
	case "query":
		fallthrough
	case "health":
		serviceHookImpl := NewServiceHookImpl(r.RestConfig)

		handlers := serviceHookImpl.GetRawAPIHandlers(serviceHookImpl.GetterFn())
		i := handlers[0].Handler(ctx, kindsys.StaticMetadata{
			Name:      name,
			Namespace: request.NamespaceValue(ctx),
		})

		// NOTE: we don't yet have error bubbling, i could be nil but streamer should handle it

		streamer.SetInputStream(i)
	default:
		// This should never happen in theory - only configured subresource APIs will trigger
		return nil, fmt.Errorf("Unknown subresource: %s", info.Subresource)
	}

	return streamer, nil
}
