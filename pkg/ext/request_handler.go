package ext

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	restclient "k8s.io/client-go/rest"
	"k8s.io/kube-openapi/pkg/spec3"
)

type requestHandler struct {
	router *mux.Router
}

func NewRequestHandler(delegateHandler http.Handler, restConfig *restclient.Config, hooks *APIServiceHooks) (*requestHandler, error) {
	router := mux.NewRouter()
	getter := makeGetter(restConfig)

	// LEVEL = Group+Verison
	var sub *mux.Router
	prefix := "/apis/googlesheets.ext.grafana.com/v1"
	for _, v := range hooks.PluginRouteHandlers {
		if v.Level == RawAPILevelGroupVersion {
			if sub == nil {
				sub = router.PathPrefix(prefix).Subrouter()
				sub.MethodNotAllowedHandler = &methodNotAllowedHandler{}
			}

			methods, err := methodsFromSpec(v.Slug, v.Spec)
			if err != nil {
				return nil, err
			}
			sub.HandleFunc(v.Slug, v.Handler).
				Methods(methods...)
		}
	}

	// LEVEL = Namespace/Tenent
	sub = nil
	prefix += "/namespaces/{ns}"
	for _, v := range hooks.PluginRouteHandlers {
		if v.Level == RawAPILevelNamespace {
			if sub == nil {
				sub = router.PathPrefix(prefix).Subrouter()
				sub.MethodNotAllowedHandler = &methodNotAllowedHandler{}
			}

			methods, err := methodsFromSpec(v.Slug, v.Spec)
			if err != nil {
				return nil, err
			}
			sub.HandleFunc(v.Slug, v.Handler).
				Methods(methods...)

		}
	}

	// LEVEL = Resource
	sub = nil
	prefix += "/datasources/{name}"
	for _, v := range hooks.PluginRouteHandlers {
		if v.Level == RawAPILevelResource {
			if sub == nil {
				sub = router.PathPrefix(prefix).Subrouter()
				sub.MethodNotAllowedHandler = &methodNotAllowedHandler{}
			}

			methods, err := methodsFromSpec(v.Slug, v.Spec)
			if err != nil {
				return nil, err
			}
			sub.HandleFunc(v.Slug, SubresourceHandlerWrapper(v.Handler, getter)).
				Methods(methods...)
		}
	}

	// Per Gorilla Mux issue here: https://github.com/gorilla/mux/issues/616#issuecomment-798807509
	// default handler must come last
	router.PathPrefix("/").Handler(delegateHandler)

	return &requestHandler{
		router: router,
	}, nil
}

func (h *requestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(w, req)
}

func methodsFromSpec(slug string, props *spec3.PathProps) ([]string, error) {
	if props == nil {
		return []string{"GET", "POST", "PUT", "PATCH", "DELETE"}, nil
	}

	methods := make([]string, 0)
	if props.Get != nil {
		methods = append(methods, "GET")
	}
	if props.Post != nil {
		methods = append(methods, "POST")
	}
	if props.Put != nil {
		methods = append(methods, "PUT")
	}
	if props.Patch != nil {
		methods = append(methods, "PATCH")
	}
	if props.Delete != nil {
		methods = append(methods, "DELETE")
	}

	if len(methods) == 0 {
		return nil, fmt.Errorf("Invalid OpenAPI Spec for slug=%s without any methods in PathProps", slug)
	}

	return methods, nil
}

type methodNotAllowedHandler struct {
}

func (h *methodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(405) // method not allowed
	// w.Write([]byte(fmt.Sprintf("bad method: %s", req.Method)))
}
