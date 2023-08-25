package ext

import (
	"net/http"

	"github.com/gorilla/mux"
	restclient "k8s.io/client-go/rest"
)

type requestHandler struct {
	router *mux.Router
}

func NewRequestHandler(apiHandler http.Handler, restConfig *restclient.Config, hooks *APIServiceHooks) *requestHandler {
	router := mux.NewRouter()

	getter := makeGetter(restConfig)
	requestHandler := &requestHandler{
		router: router,
	}

	// LEVEL = Group+Verison
	var sub *mux.Router
	prefix := "/apis/googlesheets.ext.grafana.com/v1"
	for _, v := range hooks.PluginRouteHandlers {
		if v.Level == RawAPILevelGroupVersion {
			if sub == nil {
				sub = router.PathPrefix(prefix).Subrouter()
			}
			sub.HandleFunc(v.Slug, v.Handler)
			// TODO... methods from spec!
		}
	}

	// LEVEL = Namespace/Tenent
	sub = nil
	prefix += "/namespaces/{ns}"
	for _, v := range hooks.PluginRouteHandlers {
		if v.Level == RawAPILevelNamespace {
			if sub == nil {
				sub = router.PathPrefix(prefix).Subrouter()
			}
			sub.HandleFunc(v.Slug, v.Handler)
			// TODO... methods from spec!
		}
	}

	// LEVEL = Resource
	sub = nil
	prefix += "/datasources/{name}"
	for _, v := range hooks.PluginRouteHandlers {
		if v.Level == RawAPILevelResource {
			if sub == nil {
				sub = router.PathPrefix(prefix).Subrouter()
				// ??? does not do anything!!!
				sub.MethodNotAllowedHandler = &methodNotAllowedHandler{}
			}
			sub.HandleFunc(v.Slug, SubresourceHandlerWrapper(v.Handler, getter))
			// TODO... methods from spec!
		}
	}

	// dsSubrouter := router.
	// 	PathPrefix("/apis/googlesheets.ext.grafana.com/v1/namespaces/{ns}/datasources/{name}").
	// 	Subrouter()

	// dsSubrouter.
	// 	HandleFunc("/query", SubresourceHandlerWrapper(
	// 		hooks.PluginRouteHandlers[2].Handler, getter)).
	// 	Methods("POST")

	// dsSubrouter.
	// 	HandleFunc("/health", SubresourceHandlerWrapper(
	// 		hooks.PluginRouteHandlers[3].Handler, getter)).
	// 	Methods("GET")

	// dsSubrouter.
	// 	HandleFunc("/resource/{path:.*}", SubresourceHandlerWrapper(
	// 		hooks.PluginRouteHandlers[4].Handler, getter))

	// ?????? does not seem to do anything :(
	//dsSubrouter.MethodNotAllowedHandler = &methodNotAllowedHandler{}

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
