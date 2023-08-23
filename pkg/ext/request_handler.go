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
)

type requestHandler struct {
	router          *mux.Router
	restConfig      *restclient.Config
	serviceHookImpl *ServiceHookImpl
}

// TODO: the map below is redundant as we could iterate through []RawAPIHandler and check for a match against Path
// for a given subresource
var (
	subresourceToRawHandlerIndexMap = map[string]int{
		// 0 is left unpopulated here and 0th index in the RawAPIHandlers is setup
		// with a notFound handler just in case we inadvertently get 0 back from this map for mismatched subresources
		"query":                 1,
		"health":                2,
		"resource/spreadsheets": 3,
	}
)

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

func (handler *requestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler.router.ServeHTTP(w, req)
}

func (handler *requestHandler) destroy() {
	// TODO: register any on destroy code here and get the callback registered with API Server
}

func (handler *requestHandler) callResourceHandler(writer http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	info, ok := request.RequestInfoFrom(ctx)
	if !ok {
		responsewriters.ErrorNegotiated(
			apierrors.NewInternalError(fmt.Errorf("no RequestInfo found in the context")),
			Codecs, schema.GroupVersion{}, writer, req,
		)
		return
	}

	vars := mux.Vars(req)
	resourcePath, ok := vars["resourcePath"]
	subresource := fmt.Sprintf("resource/%s", resourcePath)
	if !ok {
		writer.WriteHeader(404)
		writer.Write([]byte(fmt.Sprintf("Unknown empty resource path specified for CallResource")))
		return
	}

	switch resourcePath {
	case "spreadsheets":
		handlers := handler.serviceHookImpl.GetRawAPIHandlers(handler.serviceHookImpl.GetterFn())
		handlerWithSetupCompleted, err := handlers[subresourceToRawHandlerIndexMap[subresource]].Handler(ctx, kindsys.StaticMetadata{
			Name: info.Name,
			// curious: request.NamespaceValue(ctx) doesn't seem to be set but info.Namespace is
			Namespace: info.Namespace,
		})

		if err != nil {
			klog.Errorf("Error when getting the handler that closes on StaticMetadata: %s", err)
			writer.WriteHeader(404)
			writer.Write([]byte(""))
			return
		}

		handlerWithSetupCompleted(writer, req)
		return
	default:
		writer.WriteHeader(404)
		writer.Write([]byte(fmt.Sprintf("Unknown resource path specified for CallResource: %s", resourcePath)))
	}
}

func (handler *requestHandler) subresourceHandler(writer http.ResponseWriter, req *http.Request) {
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
		handlers := handler.serviceHookImpl.GetRawAPIHandlers(handler.serviceHookImpl.GetterFn())
		handlerWithSetupCompleted, err := handlers[subresourceToRawHandlerIndexMap[info.Subresource]].Handler(ctx, kindsys.StaticMetadata{
			Name: info.Name,
			// curious: request.NamespaceValue(ctx) doesn't seem to be set but info.Namespace is
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
