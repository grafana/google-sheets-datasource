package ext

import (
	"context"
	"fmt"
	"github.com/grafana/google-sheets-datasource/pkg/apiserver/registry"
	"github.com/grafana/grafana-apiserver/pkg/storage/filepath"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"github.com/grafana/google-sheets-datasource/pkg/apiserver"
	clientRest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net"
	"net/http"
	"os"
	"path"

	"github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets/install"
	"github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets/v1"
	"github.com/grafana/grafana-apiserver/pkg/certgenerator"
	grafanaapiserveroptions "github.com/grafana/grafana-apiserver/pkg/cmd/server/options"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/options"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	install.Install(Scheme)

	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}

func Start(ctx context.Context) error {
	// logger := logr.New(newLogAdapter())
	// logger.V(9)
	// klog.SetLoggerWithOptions(logger, klog.ContextualLogger(true))

	o := grafanaapiserveroptions.NewGrafanaAPIServerOptions(os.Stdout, os.Stderr)
	o.RecommendedOptions.SecureServing.BindPort = 6443
	// o.RecommendedOptions.Authentication.DisableAnonymous = false
	o.RecommendedOptions.Authentication.RemoteKubeConfigFile = "/Users/charandas/.kube/config"
	o.RecommendedOptions.Authorization.RemoteKubeConfigFile = "/Users/charandas/.kube/config"
	o.RecommendedOptions.Authorization.AlwaysAllowPaths = []string{"*"}
	o.RecommendedOptions.Authorization.AlwaysAllowGroups = []string{user.Anonymous}
	// o.RecommendedOptions.Authorization.
	o.RecommendedOptions.Etcd = nil
	// TODO: setting CoreAPI to nil currently segfaults in grafana-apiserver
	o.RecommendedOptions.CoreAPI = nil

	// Get the util to get the paths to pre-generated certs
	certUtil := certgenerator.CertUtil{
		K8sDataPath: "data",
	}

	err := certUtil.InitializeCACertPKI()
	if err != nil {
		fmt.Println("Err", err)
		panic("could not provision certs")
	}

	err = certUtil.EnsureApiServerPKI("127.0.0.1")
	if err != nil {
		fmt.Println("Err", err)
		panic("could not provision certs")
	}

	err = certUtil.EnsureAuthzClientPKI()
	if err != nil {
		fmt.Printf("error ensuring K8s Authz Client PKI", "error", err)
		panic("could not provision certs")
	}

	err = certUtil.EnsureAuthnClientPKI()
	if err != nil {
		fmt.Printf("error ensuring K8s Authn Client PKI", "error", err)
		panic("could not provision certs")
	}

	o.RecommendedOptions.SecureServing.BindAddress = net.ParseIP(certgenerator.DefaultAPIServerIp)
	o.RecommendedOptions.SecureServing.ServerCert.CertKey = options.CertKey{
		CertFile: certUtil.APIServerCertFile(),
		KeyFile:  certUtil.APIServerKeyFile(),
	}

	if err := o.Complete(); err != nil {
		return err
	}

	if err := o.Validate(); err != nil {
		return err
	}

	serverConfig, err := o.Config()
	if err != nil {
		return err
	}

	// rootCert, err := certUtil.GetK8sCACert()
	// if err != nil {
	// 	return err
	// }

	// authenticator, err := newAuthenticator(rootCert)
	// if err != nil {
	// 	return err
	// }

	// serverConfig.GenericConfig.Authentication.Authenticator = authenticator

	delegationTarget := genericapiserver.NewEmptyDelegate()
	delegateHandler := delegationTarget.UnprotectedHandler()
	if delegateHandler == nil {
		delegateHandler = http.NotFoundHandler()
	}

	server, err := serverConfig.Complete().New(delegationTarget)
	if err != nil {
		return err
	}

	// server.GenericAPIServer.Authorizer = &FakeAuthorizer{}

	restConfig := server.GenericAPIServer.LoopbackClientConfig
	// err = s.writeKubeConfiguration(s.restConfig)
	// if err != nil {
	// 	return err
	// }

	prepared := server.GenericAPIServer.PrepareRun()
	fmt.Printf("TODO: %v, %v\n", prepared, restConfig)

	err = writeKubeConfiguration(server.GenericAPIServer.LoopbackClientConfig)
	if err != nil {
		return err
	}

	subresourceHandler := &apiserver.SubresourceHandler{
		Storage:             nil,
		Authorizer:          server.GenericAPIServer.Authorizer,
		MaxRequestBodyBytes: serverConfig.GenericConfig.MaxRequestBodyBytes,
		DelegateHandler:     delegateHandler,
	}

	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(apiserver.PluginAPIGroup, Scheme, metav1.ParameterCodec, Codecs)
	storageMap := map[string]rest.Storage{}

	datasourceREST, err := registry.NewREST(Scheme, filepath.NewRESTOptionsGetter("/tmp/plugin-apiserver", Codecs.LegacyCodec(v1.SchemeGroupVersion)))
	if err != nil {
		return err
	}
	storageMap["datasources"] = datasourceREST
	apiGroupInfo.VersionedResourcesStorageMap[v1.SchemeGroupVersion.Version] = storageMap

	// storageMap["datasources/query"] = &storage.SubresourceStreamerREST{}

	if err := server.GenericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
		fmt.Println("Could not install API Group", err)
		return err
	}

	server.GenericAPIServer.Handler.NonGoRestfulMux.Handle(fmt.Sprintf("/apis/%s", apiserver.PluginAPIGroup), subresourceHandler)
	server.GenericAPIServer.Handler.NonGoRestfulMux.HandlePrefix(fmt.Sprintf("/apis/%s/", apiserver.PluginAPIGroup), subresourceHandler)

	// s.handler = func(c *contextmodel.ReqContext) {
	// 	req := c.Req
	// 	req.URL.Path = strings.TrimPrefix(req.URL.Path, "/k8s")
	// 	if req.URL.Path == "" {
	// 		req.URL.Path = "/"
	// 	}
	// 	ctx := req.Context()
	// 	signedInUser := appcontext.MustUser(ctx)

	// 	req.Header.Set("X-Remote-User", strconv.FormatInt(signedInUser.UserID, 10))
	// 	req.Header.Set("X-Remote-Group", "grafana")
	// 	req.Header.Set("X-Remote-Extra-token-name", signedInUser.Name)
	// 	req.Header.Set("X-Remote-Extra-org-role", string(signedInUser.OrgRole))
	// 	req.Header.Set("X-Remote-Extra-org-id", strconv.FormatInt(signedInUser.OrgID, 10))
	// 	req.Header.Set("X-Remote-Extra-user-id", strconv.FormatInt(signedInUser.UserID, 10))

	// 	resp := responsewriter.WrapForHTTP1Or2(c.Resp)
	// 	prepared.GenericAPIServer.Handler.ServeHTTP(resp, req)
	// }

	fmt.Println("Potato")
	go func() {
		c := make(chan struct{})
		err := prepared.Run(c)
		if err != nil {
			fmt.Printf("Could not run", err)
		}
	}()

	return nil
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
