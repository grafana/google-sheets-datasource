package ext

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	restclient "k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	"github.com/grafana/kindsys"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/request"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

type requestHandler struct {
	router             *mux.Router
	config             *genericapiserver.Config
	subresourceHandler *subresourceHandler
}

func NewRequestHandler(apiHandler http.Handler, c *genericapiserver.Config) *requestHandler {
	router := mux.NewRouter()

	requestHandler := &requestHandler{
		router:             router,
		config:             c,
		subresourceHandler: newSubresourceHandler(c.LoopbackClientConfig),
	}

	// Not calling DefaultBuildHandlerChain on any of the handlers written by us prevents being able to get requestInfo
	// One could argue that Gorilla Mux does provide vars based matching on route params
	// But I went with using requestInfo just to make the code looking like K8s
	router.Handle("/apis/googlesheets.ext.grafana.com/v1/namespaces/{ns}/datasources/{datasourceName}/{subresource}", genericapiserver.DefaultBuildHandlerChain(requestHandler.subresourceHandler, c))
	// Per Gorilla Mux issue here: https://github.com/gorilla/mux/issues/616#issuecomment-798807509
	// default handler must come last
	router.PathPrefix("/").Handler(genericapiserver.DefaultBuildHandlerChain(apiHandler, c))

	return requestHandler
}

func (handler *requestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler.router.ServeHTTP(w, req)
}

func (handler *requestHandler) destroy() {
	// TODO: register any on destroy code here and get the callback registered with API Server
}

// Needed to make another struct for subresourceHandler so as to satisfy http.Handler interface
// and to be able to call genericapiserver.DefaultBuildHandlerChain on it - an anonymous function sadly
// can't be wrapped like that
type subresourceHandler struct {
	restConfig *restclient.Config
}

func (subhandler *subresourceHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	info, ok := request.RequestInfoFrom(ctx)
	if !ok {
		responsewriters.ErrorNegotiated(
			apierrors.NewInternalError(fmt.Errorf("no RequestInfo found in the context")),
			Codecs, schema.GroupVersion{}, writer, req,
		)
		return
	}

	switch info.Subresource {
	case "query":
		fallthrough
	case "health":
		serviceHookImpl := NewServiceHookImpl(subhandler.restConfig)

		handlers := serviceHookImpl.GetRawAPIHandlers(serviceHookImpl.GetterFn())
		handlerWithSetupCompleted, err := handlers[0].Handler(ctx, kindsys.StaticMetadata{
			Name: info.Name,
			// curious: but request.NamespaceValue(ctx) doesn't seem to be set but info.Namespace is
			Namespace: info.Namespace,
		})

		if err != nil {
			klog.Errorf("Error when getting the handler that closes on StaticMetadata: %s", err)
			writer.WriteHeader(404)
			writer.Write([]byte(""))
			return
		}

		handlerWithSetupCompleted(writer, req)
	default:
		// This should never happen in theory - only configured subresource APIs will trigger
		writer.WriteHeader(404)
		writer.Write([]byte(fmt.Sprintf("Unknown subresource: %s", info.Subresource)))
	}
}

func newSubresourceHandler(restConfig *restclient.Config) *subresourceHandler {
	return &subresourceHandler{
		restConfig: restConfig,
	}
}
