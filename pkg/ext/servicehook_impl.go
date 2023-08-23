package ext

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	v1 "github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets/v1"
	"github.com/grafana/google-sheets-datasource/pkg/client/clientset/clientset"
	"github.com/grafana/google-sheets-datasource/pkg/googlesheets"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/kindsys"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

type InternalHandler = func(ctx context.Context, pluginCtx *backend.PluginContext, datasource *googlesheets.Datasource, path ...string) func(w http.ResponseWriter, req *http.Request)

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
				Path:    "not-found-handler",
				OpenAPI: "",
				Level:   RawAPILevel(RawAPILevelResource),
				Handler: func(_ context.Context, _ kindsys.StaticMetadata) (http.HandlerFunc, error) {
					return http.NotFound, nil
				},
			},
			{
				Path:    "query",
				OpenAPI: "",
				Level:   RawAPILevel(RawAPILevelResource),
				Handler: setupPluginContextAndReturnHandler(getter, getQueryHandler),
			},
			{
				Path:    "health",
				OpenAPI: "",
				Level:   RawAPILevel(RawAPILevelResource),
				Handler: setupPluginContextAndReturnHandler(getter, getHealthHandler),
			},
			{
				Path:    "resource/spreadsheets",
				OpenAPI: "",
				Level:   RawAPILevel(RawAPILevelResource),
				Handler: setupPluginContextAndReturnHandler(getter, getCallResourceHandler, "/spreadsheets"),
			},
		}
	}
}

func setupPluginContextAndReturnHandler(getter ResourceGetter, internalHandler InternalHandler, path ...string) ClosedOnFetchedK8sResourceHandler {
	return func(ctx context.Context, id kindsys.StaticMetadata) (http.HandlerFunc, error) {
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

		// k8s HACK! move from spec to decrypted
		settings.DecryptedSecureJSONData["privateKey"] = ds.Spec.PrivateKey

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

		return internalHandler(ctx, pluginCtx, googleSheetDatasource, path...), nil
	}
}

func getCallResourceHandler(ctx context.Context, pluginCtx *backend.PluginContext, datasource *googlesheets.Datasource, path ...string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			klog.Errorf("CallResourceRequest body was malformed: %s", err)
			w.WriteHeader(400)
			w.Write([]byte("CallResourceRequest body was malformed"))
			return
		}

		wrappedSender := callResourceResponseSenderFunc(func(response *backend.CallResourceResponse) error {
			w.WriteHeader(response.Status)
			for key, headerValues := range response.Headers {
				for _, value := range headerValues {
					w.Header().Set(key, value)
				}
			}
			w.Write(response.Body)
			return nil
		})

		err = datasource.CallResource(ctx, &backend.CallResourceRequest{
			PluginContext: *pluginCtx,
			Path:          path[0],
			Method:        req.Method,
			Body:          body,
		}, wrappedSender)

		if err != nil {
			// our wrappedSender func will likely never be invoked for errors
			// respond with a 400
			w.WriteHeader(400)
			w.Write([]byte("encountered error invoking CallResponseHandler for request"))
		}
	}
}

func getHealthHandler(ctx context.Context, pluginCtx *backend.PluginContext, datasource *googlesheets.Datasource, _ ...string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		healthResponse, err := datasource.CheckHealth(ctx, &backend.CheckHealthRequest{
			PluginContext: *pluginCtx,
		})

		if err != nil {
			// our wrappedSender func will likely never be invoked for errors
			// respond with a 400
			w.WriteHeader(400)
			klog.Errorf("encountered error invoking CheckHealth: %s", err)
			w.Write([]byte("encountered error invoking CheckHealth"))
		}

		jsonRsp, err := json.Marshal(healthResponse)
		if err != nil {
			return
		}
		w.WriteHeader(200)
		w.Write(jsonRsp)
	}
}

func getQueryHandler(ctx context.Context, pluginCtx *backend.PluginContext, datasource *googlesheets.Datasource, _ ...string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			klog.Errorf("QueryDataRequest was malformed: %s", err)
			w.WriteHeader(400)
			w.Write([]byte("QueryDataRequest was malformed"))
			return
		}
		queries, err := readQueries(body)
		if err != nil {
			klog.Errorf("Could not parse QueryDataRequest: %s", err)
			w.WriteHeader(400)
			w.Write([]byte("Could not parse QueryDataRequest"))
			return
		}

		queryResponse, err := datasource.QueryData(ctx, &backend.QueryDataRequest{
			PluginContext: *pluginCtx,
			Queries:       queries,
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
	}
}

type callResourceResponseSenderFunc func(res *backend.CallResourceResponse) error

func (fn callResourceResponseSenderFunc) Send(res *backend.CallResourceResponse) error {
	return fn(res)
}
