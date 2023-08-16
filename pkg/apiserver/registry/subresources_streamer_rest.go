package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/grafana/google-sheets-datasource/pkg/googlesheets"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"k8s.io/apiserver/pkg/endpoints/request"

	v1 "github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets/v1"
	"github.com/grafana/google-sheets-datasource/pkg/apiserver/apihelpers"
	"github.com/grafana/google-sheets-datasource/pkg/client/clientset/clientset"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	restclient "k8s.io/client-go/rest"
)

var _ rest.Storage = (*SubresourceStreamerREST)(nil)
var _ rest.Getter = (*SubresourceStreamerREST)(nil)

type SubresourceStreamerREST struct {
	// store *genericregistry.Store
	RestConfig *restclient.Config
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
	return &v1.Datasource{}
}

func (r *SubresourceStreamerREST) Destroy() {
}

/* func (r *SubresourceStreamerREST) New() runtime.Object {
	return r.store.New()
} */

func (r *SubresourceStreamerREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	cs, err := clientset.NewForConfig(r.RestConfig)
	if err != nil {
		return nil, err
	}

	ds, err := cs.GooglesheetsV1().Datasources(request.NamespaceValue(ctx)).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	settings := backend.DataSourceInstanceSettings{}
	settings.JSONData, err = json.Marshal(ds.Spec)
	settings.DecryptedSecureJSONData = map[string]string{}

	settings.DecryptedSecureJSONData["apiKey"] = ds.Spec.APIKey
	settings.DecryptedSecureJSONData["jwt"] = ds.Spec.JWT

	settings.Type = "grafana-googlesheets-datasource"

	pluginCtx := &backend.PluginContext{
		OrgID:                      0,
		PluginID:                   settings.Type,
		User:                       &backend.User{},
		AppInstanceSettings:        &backend.AppInstanceSettings{},,
		DataSourceInstanceSettings: settings,
	}

	instance, err := googlesheets.NewDatasource(settings)

	if err != nil {
		return nil, err
	}

	i := func(ctx context.Context, apiVersion, acceptHeader string) (stream io.ReadCloser, flush bool, mimeType string, err error) {
		instance.QueryData(ctx, &backend.QueryDataRequest{
			PluginContext: pluginCtx,
			Queries: []backend.DataQuery{},
			//  Headers: // from context
		})

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
