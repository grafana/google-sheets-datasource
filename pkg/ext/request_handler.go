package ext

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/request"
	restclient "k8s.io/client-go/rest"
)

type requestHandler struct {
	router          *mux.Router
	restConfig      *restclient.Config
	serviceHookImpl *ServiceHookImpl
}

// restclient.Config is only used by subresource handler, so we don't save it on the resource handler
func NewRequestHandler(apiHandler http.Handler, restConfig *restclient.Config) *requestHandler {
	router := mux.NewRouter()

	requestHandler := &requestHandler{
		router:          router,
		restConfig:      restConfig,
		serviceHookImpl: NewServiceHookImpl(restConfig),
	}

	dsSubrouter := router.
		PathPrefix("/apis/googlesheets.ext.grafana.com/v1/namespaces/{ns}/datasources/{datasourceName}").
		Subrouter()

	dsSubrouter.
		HandleFunc("/query", requestHandler.subresourceHandler).
		Methods("POST")

	dsSubrouter.
		HandleFunc("/health", requestHandler.subresourceHandler).
		Methods("GET")

	dsSubrouter.
		HandleFunc("/resource/{resourcePath:.*}", requestHandler.callResourceHandler)

	// Per Gorilla Mux issue here: https://github.com/gorilla/mux/issues/616#issuecomment-798807509
	// default handler must come last
	router.PathPrefix("/").Handler(apiHandler)

	return requestHandler
}

func (h *requestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(w, req)
}

func (h *requestHandler) destroy() {
	// TODO: register any on destroy code here and get the callback registered with API Server
}

func (h *requestHandler) callResourceHandler(writer http.ResponseWriter, req *http.Request) {
	SubresourceHandlerWrapper(h.serviceHookImpl.PluginRouteHandlers[4].Handler, h.serviceHookImpl.GetterFn())(writer, req)

}

func (h *requestHandler) subresourceHandler(writer http.ResponseWriter, req *http.Request) {
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
		handler := SubresourceHandlerWrapper(h.serviceHookImpl.PluginRouteHandlers[2].Handler, h.serviceHookImpl.GetterFn())
		handler(writer, req)
	case "health":
		handler := SubresourceHandlerWrapper(h.serviceHookImpl.PluginRouteHandlers[3].Handler, h.serviceHookImpl.GetterFn())
		handler(writer, req)
	default:
		// This should never happen in theory - only configured subresource APIs will trigger
		writer.WriteHeader(404)
		writer.Write([]byte(fmt.Sprintf("Unknown subresource: %s", info.Subresource)))
	}
}
