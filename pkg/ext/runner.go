package ext

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/grafana/grafana-apiserver/pkg/certgenerator"
	grafanaapiserveroptions "github.com/grafana/grafana-apiserver/pkg/cmd/server/options"
	"k8s.io/apiserver/pkg/authentication/user"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/options"
)

func RunServer() error {
	http.HandleFunc("/", getRoot)
	fmt.Printf("starting k8s server on port 3333...\n")
	err := start(context.Background())
	if err != nil {
		return err
	}
	return http.ListenAndServe(":3333", nil)
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "hello!")
}

func start(ctx context.Context) error {
	// logger := logr.New(newLogAdapter())
	// logger.V(9)
	// klog.SetLoggerWithOptions(logger, klog.ContextualLogger(true))

	o := grafanaapiserveroptions.NewGrafanaAPIServerOptions(os.Stdout, os.Stderr)
	o.RecommendedOptions.SecureServing.BindPort = 6443
	o.RecommendedOptions.Authentication.RemoteKubeConfigFileOptional = true
	o.RecommendedOptions.Authorization.RemoteKubeConfigFileOptional = true
	o.RecommendedOptions.Authorization.AlwaysAllowPaths = []string{"*"}
	o.RecommendedOptions.Authorization.AlwaysAllowGroups = []string{user.SystemPrivilegedGroup, "grafana"}
	o.RecommendedOptions.Etcd = nil
	// TODO: setting CoreAPI to nil currently segfaults in grafana-apiserver
	o.RecommendedOptions.CoreAPI = nil

	// Get the util to get the paths to pre-generated certs
	certUtil := certgenerator.CertUtil{
		//K8sDataPath: s.dataPath,
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

	server, err := serverConfig.Complete().New(genericapiserver.NewEmptyDelegate())
	if err != nil {
		return err
	}

	restConfig := server.GenericAPIServer.LoopbackClientConfig
	// err = s.writeKubeConfiguration(s.restConfig)
	// if err != nil {
	// 	return err
	// }

	prepared := server.GenericAPIServer.PrepareRun()
	fmt.Printf("TODO: %v, %v\n", prepared, restConfig)

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

	// go func() {
	// 	s.stoppedCh <- prepared.Run(s.stopCh)
	// }()

	return nil
}
