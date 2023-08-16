package ext

import (
	"fmt"
	"github.com/grafana/grafana-apiserver/pkg/storage/filepath"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/managedfields"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/endpoints/handlers"
	"k8s.io/apiserver/pkg/registry/rest"
	"net/http"
	"strings"

	pluginRuntime "github.com/grafana/google-sheets-datasource/pkg/apiserver/runtime"
	"github.com/grafana/google-sheets-datasource/pkg/apiserver/storage"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

const PluginAPIGroup = "googlesheets.ext.grafana.com"
const PluginAPIVersion = "v1"

var _ http.Handler = &SubresourceHandler{}

type SubresourceHandler struct {
	Storage             *storage.PluginResourceStorage
	Authorizer          authorizer.Authorizer
	MaxRequestBodyBytes int64
	DelegateHandler     http.Handler
}

func (sh *SubresourceHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	requestInfo, ok := apirequest.RequestInfoFrom(ctx)
	if !ok {
		responsewriters.ErrorNegotiated(
			apierrors.NewInternalError(fmt.Errorf("no RequestInfo found in the context")),
			Codecs, schema.GroupVersion{}, w, req,
		)
		return
	}
	if !requestInfo.IsResourceRequest {
		pathParts := splitPath(requestInfo.Path)
		// only match /apis/<group>/<version>
		// only registered under /apis
		if len(pathParts) == 3 {
			// r.versionDiscoveryHandler.ServeHTTP(w, req)
			return
		}
		// only match /apis/<group>
		if len(pathParts) == 2 {
			// r.groupDiscoveryHandler.ServeHTTP(w, req)
			return
		}

		sh.DelegateHandler.ServeHTTP(w, req)
		return
	}

	if !ok {
		responsewriters.ErrorNegotiated(
			apierrors.NewInternalError(fmt.Errorf("no RequestInfo found in the context")),
			Codecs, schema.GroupVersion{}, w, req,
		)
		return
	}
	if !requestInfo.IsResourceRequest {
		w.WriteHeader(404)
	}

	// verb := strings.ToUpper(requestInfo.Verb)
	subresource := requestInfo.Subresource
	// scope := metrics.CleanScope(requestInfo)
	/* supportedTypes := []string{
		string(types.JSONPatchType),
		string(types.MergePatchType),
		string(types.ApplyPatchType),
	} */

	creator := pluginRuntime.NewObjectCreator()
	typer := pluginRuntime.NewObjectTyper()
	convertor := pluginRuntime.NewObjectConvertor()

	negotiatedSerializer := pluginRuntime.NewNegotiatedSerializer(Scheme)
	var standardSerializers []runtime.SerializerInfo
	for _, s := range negotiatedSerializer.SupportedMediaTypes() {
		if s.MediaType == runtime.ContentTypeProtobuf {
			continue
		}
		standardSerializers = append(standardSerializers, s)
	}

	resource := schema.GroupVersionResource{Group: PluginAPIGroup, Version: PluginAPIVersion, Resource: "datasources"}
	singularResource := schema.GroupVersionResource{Group: PluginAPIGroup, Version: PluginAPIVersion, Resource: "datasource"}
	kind := schema.GroupVersionKind{Group: PluginAPIGroup, Version: PluginAPIVersion, Kind: "Datasource"}

	table := rest.NewDefaultTableConvertor(kind.GroupVersion().WithResource("datasources").GroupResource())

	equivalentResourceRegistry := runtime.NewEquivalentResourceRegistry()
	equivalentResourceRegistry.RegisterKindFor(resource, "", kind)

	sh.Storage = storage.NewStorage(
		resource.GroupResource(),
		singularResource.GroupResource(),
		kind,
		schema.GroupVersionKind{Group: PluginAPIGroup, Version: PluginAPIVersion, Kind: "DatasourceList"},
		filepath.NewRESTOptionsGetter("/tmp/plugin-apiserver", unstructured.UnstructuredJSONScheme),
		table,
		*typer,
	)

	parameterScheme := runtime.NewScheme()
	parameterScheme.AddUnversionedTypes(schema.GroupVersion{Group: PluginAPIGroup, Version: PluginAPIVersion},
		&metav1.ListOptions{},
		&metav1.GetOptions{},
		&metav1.DeleteOptions{},
	)

	reqScope := handlers.RequestScope{
		Namer: handlers.ContextBasedNaming{
			Namer:         meta.NewAccessor(),
			ClusterScoped: false,
		},
		Serializer:          negotiatedSerializer,
		ParameterCodec:      runtime.NewParameterCodec(parameterScheme),
		StandardSerializers: standardSerializers,

		Creater:         creator,
		Convertor:       convertor,
		Defaulter:       Scheme,
		Typer:           typer,
		UnsafeConvertor: convertor,

		EquivalentResourceMapper: equivalentResourceRegistry,

		Resource: schema.GroupVersionResource{Group: PluginAPIGroup, Version: PluginAPIVersion, Resource: "Datasource"},
		Kind:     kind,

		// a handler for a specific group-version of a custom resource uses that version as the in-memory representation
		HubGroupVersion: kind.GroupVersion(),

		MetaGroupVersion: metav1.SchemeGroupVersion,

		TableConvertor: table,

		Authorizer: sh.Authorizer,

		MaxRequestBodyBytes: sh.MaxRequestBodyBytes,
	}

	resetFields := map[fieldpath.APIVersion]*fieldpath.Set{}
	reqScope, _ = scopeWithFieldManager( // TODO: check for error
		managedfields.NewDeducedTypeConverter(),
		reqScope,
		resetFields,
		"",
	)

	switch {
	case len(subresource) == 0:
		w.WriteHeader(404)
		// refs
		// history
		// query
	default:
		responsewriters.ErrorNegotiated(
			apierrors.NewNotFound(schema.GroupResource{Group: requestInfo.APIGroup, Resource: requestInfo.Resource}, requestInfo.Name),
			Codecs, schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}, w, req,
		)
	}
}

func scopeWithFieldManager(typeConverter managedfields.TypeConverter, reqScope handlers.RequestScope, resetFields map[fieldpath.APIVersion]*fieldpath.Set, subresource string) (handlers.RequestScope, error) {
	fieldManager, err := managedfields.NewDefaultCRDFieldManager(
		typeConverter,
		reqScope.Convertor,
		reqScope.Defaulter,
		reqScope.Creater,
		reqScope.Kind,
		reqScope.HubGroupVersion,
		subresource,
		resetFields,
	)
	if err != nil {
		return handlers.RequestScope{}, err
	}
	reqScope.FieldManager = fieldManager
	return reqScope, nil
}

// splitPath returns the segments for a URL path.
func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}
