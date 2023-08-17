package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/grafana/google-sheets-datasource/pkg/googlesheets"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/klog"

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
	// settings.DecryptedSecureJSONData = map[string]string{}

	// settings.DecryptedSecureJSONData["apiKey"] = ds.Spec.APIKey
	// settings.DecryptedSecureJSONData["jwt"] = ds.Spec.JWT

	settings.Type = "grafana-googlesheets-datasource"

	pluginCtx := &backend.PluginContext{
		OrgID:                      1,
		PluginID:                   settings.Type,
		User:                       &backend.User{},
		AppInstanceSettings:        &backend.AppInstanceSettings{},
		DataSourceInstanceSettings: &settings,
	}

	instance, err := googlesheets.NewDatasource(settings)
	if err != nil {
		return nil, err
	}

	googleSheetDatasource, ok := instance.(*googlesheets.GoogleSheetsDatasource)
	if !ok {
		return nil, err
	}

	customProperties, err := json.Marshal(map[string]interface{}{
		"cacheDurationSeconds": 300,
		"spreadsheet":          "19sbxbIdRUNOeYECMlq2D3nFwD5oVJf1m8YRHcB1UXOY",
		"range":                "To do!C6",
		"datasourceId":         4, // ds.Spec.Id
		"datasource": map[string]string{
			"uid":  "b1808c48-9fc9-4045-82d7-081781f8a553",
			"type": "grafana-googlesheets-datasource",
		},
	})

	if err != nil {
		return nil, err
	}

	i := func(ctx context.Context, apiVersion, acceptHeader string) (stream io.ReadCloser, flush bool, mimeType string, err error) {

		queryResponse, err := googleSheetDatasource.QueryData(ctx, &backend.QueryDataRequest{
			PluginContext: *pluginCtx,
			Queries: []backend.DataQuery{
				{
					RefID: "A",
					// QueryType:     "", // not defined in the original request as sniffed from a browser session
					MaxDataPoints: 1541,
					Interval:      15000,
					TimeRange:     backend.TimeRange{},
					JSON:          customProperties,
				},
			},
			//  Headers: // from context
		})
		if err != nil {
			klog.Info("QueryResponse: +%v", queryResponse)
		}

		jsonRsp, err := json.Marshal(queryResponse)
		if err != nil {
			return nil, false, "", err
		}
		return io.NopCloser(bytes.NewReader(jsonRsp)), false, "application/json", nil
	}

	streamer := &apihelpers.SubresourceStreamer{}
	streamer.SetInputStream(i)

	return streamer, nil
}
