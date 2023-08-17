package ext

import (
	"context"

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

type SubresourceStreamerREST struct {
	RestConfig *restclient.Config
}

func (r *SubresourceStreamerREST) New() runtime.Object {
	return &v1.Datasource{}
}

func (r *SubresourceStreamerREST) Destroy() {
}

func (r *SubresourceStreamerREST) Get(ctx context.Context, name string, _ *metav1.GetOptions) (runtime.Object, error) {
	serviceHookImpl := NewServiceHookImpl(r.RestConfig)

	handlers := serviceHookImpl.GetRawAPIHandlers(serviceHookImpl.GetterFn())
	i := handlers[0].Handler(ctx, kindsys.StaticMetadata{
		Name:      name,
		Namespace: request.NamespaceValue(ctx),
	})

	// NOTE: we don't yet have error bubbling, i could be nil but streamer should handle it

	streamer := &apihelpers.SubresourceStreamer{}
	streamer.SetInputStream(i)

	return streamer, nil
}
