package ext

import (
	"fmt"
	"net/http"
	"strings"
)

type RequestHandler struct {
	delegate http.Handler
}

func (r *RequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// ctx := req.Context()
	/* requestInfo, ok := apirequest.RequestInfoFrom(ctx)
	if !ok {
		responsewriters.ErrorNegotiated(
			apierrors.NewInternalError(fmt.Errorf("no RequestInfo found in the context")),
			Codecs, schema.GroupVersion{}, w, req,
		)
		return
	} */

	fmt.Println("Hello, path", req.URL.Path)
	if strings.HasSuffix(req.URL.Path, "/query") {
		w.Write([]byte("Hello!"))
		return
	}

	// only match /apis/<group>/<version>
	// only registered under /apis
	/* if len(pathParts) == 3 {
		r.versionDiscoveryHandler.ServeHTTP(w, req)
		return
	} */
	// only match /apis/<group>
	/* if len(pathParts) == 2 {
		r.groupDiscoveryHandler.ServeHTTP(w, req)
		return
	} */

	r.delegate.ServeHTTP(w, req)

}

func (r *RequestHandler) destroy() {

}
