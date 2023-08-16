package ext

import (
	"fmt"
	"github.com/grafana/google-sheets-datasource/pkg/apiserver/registry"
	"github.com/grafana/grafana-apiserver/pkg/storage/filepath"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/version"
	clientRest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"path"

	"github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets/install"
	"github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets/v1"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)

	// if you modify this, make sure you update the crEncoder
	unversionedVersion = schema.GroupVersion{Group: "", Version: "v1"}
	unversionedTypes   = []runtime.Object{
		&metav1.Status{},
		&metav1.WatchEvent{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	}
)

func init() {
	fmt.Println("init: Potato")
	install.Install(Scheme)

	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Group: "", Version: "v1"})
	Scheme.AddUnversionedTypes(unversionedVersion, unversionedTypes...)

}

type PluginAggregatedServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
}

// CompletedConfig embeds a private pointer that cannot be instantiated outside of this package.
type CompletedConfig struct {
	*completedConfig
}

func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
	}

	c.GenericConfig.Version = &version.Info{
		Major: "1",
		Minor: "0",
	}

	return CompletedConfig{&c}
}

// New returns a new instance of PluginAggregatedServer from the given config.
func (c completedConfig) New() (*PluginAggregatedServer, error) {
	genericServer, err := c.GenericConfig.New("sample-apiserver", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	s := &PluginAggregatedServer{
		GenericAPIServer: genericServer,
	}

	err = writeKubeConfiguration(s.GenericAPIServer.LoopbackClientConfig)
	if err != nil {
		return nil, err
	}

	delegationTarget := genericapiserver.NewEmptyDelegate()
	delegateHandler := delegationTarget.UnprotectedHandler()
	if delegateHandler == nil {
		delegateHandler = http.NotFoundHandler()
	}

	/* subresourceHandler := &SubresourceHandler{
		Storage:             nil,
		Authorizer:          s.GenericAPIServer.Authorizer,
		MaxRequestBodyBytes: c.GenericConfig.MaxRequestBodyBytes,
		DelegateHandler:     delegateHandler,
	} */

	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(PluginAPIGroup, Scheme, metav1.ParameterCodec, Codecs)
	storageMap := map[string]rest.Storage{}

	datasourceREST, err := registry.NewREST(Scheme, filepath.NewRESTOptionsGetter("/tmp/plugin-apiserver", Codecs.LegacyCodec(v1.SchemeGroupVersion)))
	if err != nil {
		return nil, err
	}
	storageMap["datasources"] = datasourceREST
	// storageMap["datasources/query"] = &storage.SubresourceStreamerREST{}
	apiGroupInfo.VersionedResourcesStorageMap["v1"] = storageMap

	if err := s.GenericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
		fmt.Println("Could not install API Group", err)
		return nil, err
	}

	// s.GenericAPIServer.Handler.NonGoRestfulMux.Handle(fmt.Sprintf("/apis/%s", PluginAPIGroup), subresourceHandler)
	// s.GenericAPIServer.Handler.NonGoRestfulMux.HandlePrefix(fmt.Sprintf("/apis/%s/", PluginAPIGroup), subresourceHandler)

	return s, nil
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

func writeKubeConfiguration(restConfig *clientRest.Config) error {
	clusters := make(map[string]*clientcmdapi.Cluster)
	clusters["default-cluster"] = &clientcmdapi.Cluster{
		Server:                restConfig.Host,
		InsecureSkipTLSVerify: true,
	}

	contexts := make(map[string]*clientcmdapi.Context)
	contexts["default-context"] = &clientcmdapi.Context{
		Cluster:   "default-cluster",
		Namespace: "default",
		AuthInfo:  "default",
	}

	authinfos := make(map[string]*clientcmdapi.AuthInfo)
	authinfos["default"] = &clientcmdapi.AuthInfo{
		Token: restConfig.BearerToken,
	}

	clientConfig := clientcmdapi.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       clusters,
		Contexts:       contexts,
		CurrentContext: "default-context",
		AuthInfos:      authinfos,
	}
	return clientcmd.WriteToFile(clientConfig, path.Join("data", "grafana.kubeconfig"))
}
