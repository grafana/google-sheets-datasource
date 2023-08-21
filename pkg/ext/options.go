package ext

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

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

	serverConfig := genericapiserver.NewRecommendedConfig(Codecs)
	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(openapi.GetOpenAPIDefinitionsWithoutDisabledFeatures(generatedopenapi.GetOpenAPIDefinitions), openapinamer.NewDefinitionNamer(Scheme, scheme.Scheme))
	serverConfig.OpenAPIV3Config = genericapiserver.DefaultOpenAPIV3Config(openapi.GetOpenAPIDefinitionsWithoutDisabledFeatures(generatedopenapi.GetOpenAPIDefinitions), openapinamer.NewDefinitionNamer(Scheme, scheme.Scheme))
	serverConfig.SkipOpenAPIInstallation = false
	serverConfig.SharedInformerFactory = clientGoInformers.NewSharedInformerFactory(fake.NewSimpleClientset(), 10*time.Minute)
	serverConfig.ClientConfig = &rest.Config{}
	serverConfig.BuildHandlerChainFunc = func(apiHandler http.Handler, c *genericapiserver.Config) http.Handler {
		// Not calling DefaultBuildHandlerChain on any of the handlers written by us prevents being able to get requestInfo
		// One could argue that Gorilla Mux does provide vars based matching on route params
		// But I went with using requestInfo just to make the code looking like K8s
		return genericapiserver.DefaultBuildHandlerChain(NewRequestHandler(apiHandler, c.LoopbackClientConfig), c)
	}

	if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	config := &Config{
		GenericConfig: serverConfig,
	}
	return config, nil
}
