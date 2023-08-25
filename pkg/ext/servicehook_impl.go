package ext

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	v1 "github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets/v1"
	"github.com/grafana/google-sheets-datasource/pkg/client/clientset/clientset"
	"github.com/grafana/google-sheets-datasource/pkg/googlesheets"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/kindsys"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"k8s.io/kube-openapi/pkg/spec3"
)

func NewServiceHooks() *APIServiceHooks {
	// right now this makes a new datasource for each request!
	// ideally can use instancemanager with a key on resource version!
	dsgetter := func(ctx context.Context) (*googlesheets.Datasource, *backend.PluginContext, error) {
		res, err := ResourceFromContext(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("no resource in context: %w", err)
		}

		ds, ok := res.(*v1.ResourceV0)
		if !ok {
			return nil, nil, fmt.Errorf("expected datasource resource config [%t]", res)
		}

		settings := backend.DataSourceInstanceSettings{}
		settings.JSONData, err = json.Marshal(ds.Spec)
		if err != nil {
			return nil, nil, err
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
			return nil, nil, err
		}

		googleSheetDatasource, ok := instance.(*googlesheets.Datasource)
		if !ok {
			return nil, nil, fmt.Errorf("expected googlesheets.Datasource [%t]", instance)
		}
		return googleSheetDatasource, pluginCtx, nil
	}

	return &APIServiceHooks{
		Kind:         nil,
		BeforeAdd:    nil,
		BeforeUpdate: nil,

		PluginRouteHandlers: []PluginRouteHandler{
			{
				Level: RawAPILevelGroupVersion,
				Slug:  "/xxx", // URL will be
				Spec: spec3.PathProps{
					Summary:     "an example at the root level",
					Description: "longer description here?",
					Get: &spec3.Operation{
						OperationProps: spec3.OperationProps{
							Parameters: []*spec3.Parameter{
								{ParameterProps: spec3.ParameterProps{
									Name: "xxx",
								}},
							},
						},
					},
				},
				Handler: func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("Root level handler (xxx)"))
				},
			},
			{
				Level: RawAPILevelNamespace,
				Slug:  "/yyy", // URL will be
				Spec: spec3.PathProps{
					Summary:     "an example at the namespace level",
					Description: "longer description here?",
					Get: &spec3.Operation{
						OperationProps: spec3.OperationProps{
							Parameters: []*spec3.Parameter{
								{ParameterProps: spec3.ParameterProps{
									Name: "yyy",
								}},
							},
						},
					},
				},
				Handler: func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("namespace level handler (yyy)"))
				},
			},
			{
				Level: RawAPILevelResource,
				Slug:  "/query", // URL will be
				Spec: spec3.PathProps{
					Summary:     "an example at the root level",
					Description: "longer description here?",
					Post: &spec3.Operation{
						OperationProps: spec3.OperationProps{
							RequestBody: &spec3.RequestBody{
								RequestBodyProps: spec3.RequestBodyProps{
									Required: true,
								},
							},
						},
					},
				},
				Handler: func(w http.ResponseWriter, r *http.Request) {
					ctx := r.Context()
					ds, pluginCtx, err := dsgetter(ctx)
					if err != nil {
						klog.Errorf("CallResourceRequest body was malformed: %s", err)
						w.Write([]byte("error"))
						w.WriteHeader(500)
						return
					}

					executeQueryHandler(ctx, w, r, pluginCtx, ds)
				},
			},
			{
				Level: RawAPILevelResource,
				Slug:  "/health", // URL will be
				Spec: spec3.PathProps{
					Summary:     "check if the service is OK",
					Description: "longer description here?",
					Get:         &spec3.Operation{},
				},
				Handler: func(w http.ResponseWriter, r *http.Request) {
					ctx := r.Context()
					ds, pluginCtx, err := dsgetter(ctx)
					if err != nil {
						klog.Errorf("CallResourceRequest body was malformed: %s", err)
						w.Write([]byte("error"))
						w.WriteHeader(500)
						return
					}
					executeHealthHandler(ctx, w, r, pluginCtx, ds)
				},
			},
			{
				Level: RawAPILevelResource,
				Slug:  "/resource/.*", // URL will be
				Handler: func(w http.ResponseWriter, r *http.Request) {
					ctx := r.Context()
					ds, pluginCtx, err := dsgetter(ctx)
					if err != nil {
						klog.Errorf("CallResourceRequest body was malformed: %s", err)
						w.Write([]byte("error"))
						w.WriteHeader(500)
						return
					}
					executeCallResourceHandler(ctx, w, r, pluginCtx, ds)
				},
			},
		},
	}
}

func makeGetter(restConfig *restclient.Config) ResourceGetter {
	return func(ctx context.Context, ns string, name string) (kindsys.Resource, error) {
		// TODO: until we have a real resource getter, doing that inside the hook
		cs, err := clientset.NewForConfig(restConfig)
		if err != nil {
			// log error
			return nil, err
		}

		ds, err := cs.GooglesheetsV1().Datasources(ns).Get(ctx, name, metav1.GetOptions{})
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

func executeCallResourceHandler(ctx context.Context, w http.ResponseWriter, req *http.Request, pluginCtx *backend.PluginContext, datasource *googlesheets.Datasource) {
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

	idx := strings.LastIndex(req.URL.Path, "/resource")
	if idx < 0 {
		w.WriteHeader(400)
		w.Write([]byte("expected resource path"))
		return
	}
	path := req.URL.Path[idx+len("/resource"):]

	err = datasource.CallResource(ctx, &backend.CallResourceRequest{
		PluginContext: *pluginCtx,
		Path:          path,
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

func executeHealthHandler(ctx context.Context, w http.ResponseWriter, _ *http.Request, pluginCtx *backend.PluginContext, datasource *googlesheets.Datasource, _ ...string) {

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

func executeQueryHandler(ctx context.Context, w http.ResponseWriter, req *http.Request, pluginCtx *backend.PluginContext, datasource *googlesheets.Datasource, _ ...string) {
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

type callResourceResponseSenderFunc func(res *backend.CallResourceResponse) error

func (fn callResourceResponseSenderFunc) Send(res *backend.CallResourceResponse) error {
	return fn(res)
}
