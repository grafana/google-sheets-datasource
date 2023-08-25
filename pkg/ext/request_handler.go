package ext

import (
	"net/http"

	"github.com/gorilla/mux"
	restclient "k8s.io/client-go/rest"
)

type requestHandler struct {
	router *mux.Router
}

// restclient.Config is only used by subresource handler, so we don't save it on the resource handler
func NewRequestHandler(apiHandler http.Handler, restConfig *restclient.Config) *requestHandler {
	router := mux.NewRouter()

	getter := makeGetter(restConfig)
	hooks := NewServiceHooks()
	requestHandler := &requestHandler{
		router: router,
	}

	dsSubrouter := router.
		PathPrefix("/apis/googlesheets.ext.grafana.com/v1/namespaces/{ns}/datasources/{name}").
		Subrouter()

	dsSubrouter.
		HandleFunc("/query", SubresourceHandlerWrapper(
			hooks.PluginRouteHandlers[2].Handler, getter)).
		Methods("POST")

	dsSubrouter.
		HandleFunc("/health", SubresourceHandlerWrapper(
			hooks.PluginRouteHandlers[3].Handler, getter)).
		Methods("GET")

	dsSubrouter.
		HandleFunc("/resource/{path:.*}", SubresourceHandlerWrapper(
			hooks.PluginRouteHandlers[4].Handler, getter))

	// ?????? does not seem to do anything :(
	dsSubrouter.MethodNotAllowedHandler = &methodNotAllowedHandler{}

	// Per Gorilla Mux issue here: https://github.com/gorilla/mux/issues/616#issuecomment-798807509
	// default handler must come last
	router.PathPrefix("/").Handler(apiHandler)

	return requestHandler
}

func (h *requestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(w, req)
}

type methodNotAllowedHandler struct {
}

func (h *methodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(405) // method not allowed
	// w.Write([]byte(fmt.Sprintf("bad method: %s", req.Method)))
}
