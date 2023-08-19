package ext

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	v1 "github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets/v1"
	"github.com/grafana/google-sheets-datasource/pkg/client/clientset/clientset"
	"github.com/grafana/google-sheets-datasource/pkg/googlesheets"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	restclient "k8s.io/client-go/rest"
)

import "github.com/grafana/kindsys"

type ServiceHookImpl struct {
	*APIServiceHooks
	RestConfig *restclient.Config
}

func NewServiceHookImpl(restConfig *restclient.Config) *ServiceHookImpl {
	shi := &ServiceHookImpl{
		APIServiceHooks: &APIServiceHooks{
			Kind:              nil,
			BeforeAdd:         nil,
			BeforeUpdate:      nil,
			GetRawAPIHandlers: nil,
		},
		RestConfig: restConfig,
	}
	shi.setupGetRawAPIHandlers()
	return shi
}

func (shi *ServiceHookImpl) GetterFn() ResourceGetter {
	return func(ctx context.Context, id kindsys.StaticMetadata) (kindsys.Resource, error) {
		// TODO: until we have a real resource getter, doing that inside the hook
		cs, err := clientset.NewForConfig(shi.RestConfig)
		if err != nil {
			// log error
			return nil, err
		}

		ds, err := cs.GooglesheetsV1().Datasources(id.Namespace).Get(ctx, id.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		r := &v1.ResourceV0{}
		r.Spec = ds.Spec
		r.Status = ds.Status
		r.SetCommonMetadata(kindsys.CommonMetadata{
			UID:             string(ds.UID),
			ResourceVersion: ds.ResourceVersion,
			Labels:          ds.Labels,
			Finalizers:      ds.Finalizers,
			CreatedBy:       "",
			UpdatedBy:       "",
			Origin:          nil,
			ExtraFields:     nil,
		})

		r.SetStaticMetadata(kindsys.StaticMetadata{
			Group:     ds.GroupVersionKind().Group,
			Version:   ds.GroupVersionKind().Version,
			Kind:      ds.GroupVersionKind().Kind,
			Namespace: ds.ObjectMeta.Namespace,
			Name:      ds.ObjectMeta.Name,
		})

		return r, nil
	}
}

func (shi *ServiceHookImpl) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}

func (shi *ServiceHookImpl) setupGetRawAPIHandlers() {
	shi.GetRawAPIHandlers = func(getter ResourceGetter) []RawAPIHandler {
		// TODO: until we have a real resource getter, doing that inside the hook
		return []RawAPIHandler{
			{
				Path:    "/{subresource}/query",
				OpenAPI: "",
				Level:   RawAPILevel(RawAPILevelResource),
				Handler: func(ctx context.Context, id kindsys.StaticMetadata) (http.HandlerFunc, error) {
					r, err := getter(ctx, id)
					if err != nil {
						return nil, err
					}

					ds, ok := r.(*v1.ResourceV0)
					if !ok {
						return nil, fmt.Errorf("type assertion failed to kindsys ResourceV0")
					}

					settings := backend.DataSourceInstanceSettings{}
					settings.JSONData, err = json.Marshal(ds.Spec)

					if err != nil {
						return nil, err
					}

					settings.DecryptedSecureJSONData = map[string]string{}
					settings.DecryptedSecureJSONData["apiKey"] = ds.Spec.APIKey
					settings.DecryptedSecureJSONData["jwt"] = ds.Spec.JWT

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

					googleSheetDatasource, ok := instance.(*googlesheets.Datasource)
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

					return func(w http.ResponseWriter, req *http.Request) {
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
							return
						}

						jsonRsp, err := json.Marshal(queryResponse)
						if err != nil {
							return
						}
						w.WriteHeader(200)
						w.Write(jsonRsp)
					}, nil
				},
			},
		}
	}
}
