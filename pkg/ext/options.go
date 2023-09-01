package ext

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"k8s.io/kube-openapi/pkg/spec3"
	"k8s.io/kube-openapi/pkg/validation/spec"

	v1 "github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets/v1"
	"github.com/grafana/google-sheets-datasource/pkg/client/clientset/clientset"
	"github.com/grafana/google-sheets-datasource/pkg/client/clientset/clientset/scheme"
	informers "github.com/grafana/google-sheets-datasource/pkg/client/informers/externalversions"
	generatedopenapi "github.com/grafana/google-sheets-datasource/pkg/client/openapi"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/admission"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/apiserver/pkg/util/openapi"
	clientGoInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	netutils "k8s.io/utils/net"
)

type PluginAggregatedServerOptions struct {
	RecommendedOptions *genericoptions.RecommendedOptions

	SharedInformerFactory informers.SharedInformerFactory
	StdOut                io.Writer
	StdErr                io.Writer

	AlternateDNS []string
}

func NewPluginAggregatedServerOptions(out, errOut io.Writer) *PluginAggregatedServerOptions {
	o := &PluginAggregatedServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(
			"",
			Codecs.LegacyCodec(v1.SchemeGroupVersion),
		),
		StdOut: out,
		StdErr: errOut,
	}
	o.RecommendedOptions.Etcd.StorageConfig.EncodeVersioner = runtime.NewMultiGroupVersioner(v1.SchemeGroupVersion, schema.GroupKind{Group: v1.GroupName})
	return o
}

// Complete fills in fields required to have valid data
func (o *PluginAggregatedServerOptions) Complete() error {
	return nil
}

func (o *PluginAggregatedServerOptions) Config() (*Config, error) {
	// TODO have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", o.AlternateDNS, []net.IP{netutils.ParseIPSloppy("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	o.RecommendedOptions.ExtraAdmissionInitializers = func(c *genericapiserver.RecommendedConfig) ([]admission.PluginInitializer, error) {
		client, err := clientset.NewForConfig(c.LoopbackClientConfig)
		if err != nil {
			return nil, err
		}
		informerFactory := informers.NewSharedInformerFactory(client, c.LoopbackClientConfig.Timeout)
		o.SharedInformerFactory = informerFactory
		return []admission.PluginInitializer{}, nil
	}

	o.RecommendedOptions.Admission.RecommendedPluginOrder = []string{}
	o.RecommendedOptions.Admission.DisablePlugins = []string{}
	o.RecommendedOptions.Admission.EnablePlugins = []string{}

	o.RecommendedOptions.SecureServing.BindPort = 6443
	// o.RecommendedOptions.Authentication.DisableAnonymous = false
	o.RecommendedOptions.Authentication.RemoteKubeConfigFileOptional = true
	// Setting authorization to nil sets authorization to always allow more effectively than below options
	// better for development / testing
	o.RecommendedOptions.Authorization = nil
	// o.RecommendedOptions.Authorization.RemoteKubeConfigFileOptional = true
	// o.RecommendedOptions.Authorization.AlwaysAllowPaths = []string{"*"}
	// o.RecommendedOptions.Authorization.AlwaysAllowGroups = []string{user.AllUnauthenticated, user.AllAuthenticated}
	o.RecommendedOptions.Etcd.StorageConfig.Transport.ServerList = []string{"127.0.0.1:2379"}
	o.RecommendedOptions.CoreAPI = nil

	// o.RecommendedOptions.SecureServing.

	// The custom hooks!
	hooks := NewServiceHooks()

	serverConfig := genericapiserver.NewRecommendedConfig(Codecs)
	serverConfig.CorsAllowedOriginList = []string{".*"}
	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(openapi.GetOpenAPIDefinitionsWithoutDisabledFeatures(generatedopenapi.GetOpenAPIDefinitions), openapinamer.NewDefinitionNamer(Scheme, scheme.Scheme))
	serverConfig.OpenAPIConfig.PostProcessSpec = func(s *spec.Swagger) (*spec.Swagger, error) {
		s.Info.Title = "POST PROCESSED!!!"
		s.Info.VendorExtensible = spec.VendorExtensible{
			Extensions: map[string]any{"hello": "world"},
		}
		s.Paths.Paths["/apis/googlesheets.ext.grafana.com/v1/aaaa-custom-root-api"] = spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Get: &spec.Operation{
					OperationProps: spec.OperationProps{
						Description: "hello just a custom route",
						Produces: []string{
							"application/json",
						},
						Responses: &spec.Responses{
							ResponsesProps: spec.ResponsesProps{
								StatusCodeResponses: map[int]spec.Response{
									200: {},
									500: {},
								},
							},
						},
					},
				},
			},
		}
		s.Paths.Paths["/apis/googlesheets.ext.grafana.com/v1/namespaces/{namespace}/aaaa-custom-tenant-level-api"] = spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Get: &spec.Operation{
					OperationProps: spec.OperationProps{
						Description: "note! the route can not match a resource name (datasource)",
						Produces: []string{
							"application/json",
						},
						Responses: &spec.Responses{
							ResponsesProps: spec.ResponsesProps{
								StatusCodeResponses: map[int]spec.Response{
									200: {},
									500: {},
								},
							},
						},
					},
				},
			},
		}
		s.Paths.Paths["/apis/googlesheets.ext.grafana.com/v1/namespaces/{namespace}/datasources/{name}/query"] = spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Post: &spec.Operation{
					OperationProps: spec.OperationProps{
						Description: "The query method (currently sent to /ds/query)",
						Produces: []string{
							"application/json",
						},
						Responses: &spec.Responses{
							ResponsesProps: spec.ResponsesProps{
								StatusCodeResponses: map[int]spec.Response{
									200: {},
									500: {},
								},
							},
						},
					},
				},
			},
		}
		s.Paths.Paths["/apis/googlesheets.ext.grafana.com/v1/namespaces/{namespace}/datasources/{name}/health"] = spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Get: &spec.Operation{
					OperationProps: spec.OperationProps{
						Description: "Checks if the datasource config is OK",
						Produces: []string{
							"application/json",
						},
						Responses: &spec.Responses{
							ResponsesProps: spec.ResponsesProps{
								StatusCodeResponses: map[int]spec.Response{
									200: {},
									500: {},
								},
							},
						},
					},
				},
			},
		}
		return s, nil
	}

	serverConfig.OpenAPIV3Config = genericapiserver.DefaultOpenAPIV3Config(openapi.GetOpenAPIDefinitionsWithoutDisabledFeatures(generatedopenapi.GetOpenAPIDefinitions), openapinamer.NewDefinitionNamer(Scheme, scheme.Scheme))
	serverConfig.OpenAPIV3Config.PostProcessSpec3 = func(s *spec3.OpenAPI) (*spec3.OpenAPI, error) {
		prefix := "/apis/googlesheets.ext.grafana.com/v1"
		if s.Paths != nil && s.Paths.Paths[prefix+"/"] != nil {
			copy := *s // will copy the rest of the properties
			copy.Info.Title = "Google sheets plugin!"
			for _, v := range hooks.PluginRouteHandlers {
				path := prefix
				switch v.Level {
				case RawAPILevelResource:
					path += "/namespaces/{ns}/datasources/{name}"
				case RawAPILevelNamespace:
					path += "/namespaces/{ns}"
				}
				path += v.Slug
				if v.Spec != nil {
					copy.Paths.Paths[path] = &spec3.Path{PathProps: *v.Spec}
				}
			}
			return &copy, nil
		}
		return s, nil
	}
	serverConfig.SkipOpenAPIInstallation = false
	serverConfig.SharedInformerFactory = clientGoInformers.NewSharedInformerFactory(fake.NewSimpleClientset(), 10*time.Minute)
	serverConfig.ClientConfig = &rest.Config{}
	serverConfig.BuildHandlerChainFunc = func(delegateHandler http.Handler, c *genericapiserver.Config) http.Handler {
		// Call DefaultBuildHandlerChain on the main entrypoint http.Handler
		// See https://github.com/kubernetes/apiserver/blob/v0.28.0/pkg/server/config.go#L906
		// DefaultBuildHandlerChain provides many things, notably CORS, HSTS, cache-control, authz and latency tracking
		requestHandler, err := NewRequestHandler(
			delegateHandler,
			c.LoopbackClientConfig,
			hooks)
		if err != nil {
			panic(fmt.Sprintf("could not build handler chain func: %s", err.Error()))
		}
		return genericapiserver.DefaultBuildHandlerChain(requestHandler, c)
	}

	if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	config := &Config{
		GenericConfig: serverConfig,
	}
	return config, nil
}

// Run starts a new WardleServer given PluginAggregatedServerOptions
func (o PluginAggregatedServerOptions) Run(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	server.GenericAPIServer.AddPostStartHookOrDie("start-sample-server-informers", func(context genericapiserver.PostStartHookContext) error {
		config.GenericConfig.SharedInformerFactory.Start(context.StopCh)
		o.SharedInformerFactory.Start(context.StopCh)
		return nil
	})

	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}
