package storage

import (
	"bytes"
	"context"
	"github.com/grafana/google-sheets-datasource/pkg/apiserver/apihelpers"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
)

var _ rest.Storage = (*SubresourceStreamerREST)(nil)
var _ rest.Getter = (*SubresourceStreamerREST)(nil)

type SubresourceStreamerREST struct {
	// store *genericregistry.Store
}

/* func NewSubresourceStreamerREST(resource schema.GroupResource, singularResource schema.GroupResource, strategy apihelpers.StreamerStrategy, optsGetter generic.RESTOptionsGetter, tableConvertor rest.TableConvertor) *SubresourceStreamerREST {
	var storage SubresourceStreamerREST
	store := &genericregistry.Store{
		NewFunc:     func() runtime.Object { return &apihelpers.SubresourceStreamer{} },
		NewListFunc: func() runtime.Object { return &apihelpers.SubresourceStreamer{} },

		DefaultQualifiedResource:  resource,
		SingularQualifiedResource: singularResource,

		CreateStrategy:      strategy,
		UpdateStrategy:      strategy,
		DeleteStrategy:      strategy,
		ResetFieldsStrategy: strategy,

		TableConvertor: tableConvertor,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		panic(err) // TODO: Propagate error up
	}
	storage.store = store
	return &storage

} */

func (r *SubresourceStreamerREST) New() runtime.Object {
	return &apihelpers.SubresourceStreamer{}
}

func (r *SubresourceStreamerREST) Destroy() {
}

/* func (r *SubresourceStreamerREST) New() runtime.Object {
	return r.store.New()
} */

func (r *SubresourceStreamerREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	i := func(ctx context.Context, apiVersion, acceptHeader string) (stream io.ReadCloser, flush bool, mimeType string, err error) {
		jsonRsp := []byte("{\"test\": \"true\"}")
		if err != nil {
			return nil, false, "", err
		}
		return io.NopCloser(bytes.NewReader(jsonRsp)), false, "application/json", nil
	}

	streamer := &apihelpers.SubresourceStreamer{}
	streamer.SetInputStream(i)

	return streamer, nil
}
