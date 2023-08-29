package registry

import (
	"context"

	"github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
)

// NOTE: quick and dirty method of implementing a readonly resource:
// Satisfy the required interfaces (not all of genericregistry.Store) as listed below.
// In practice, we won't be wrapping a Store inside a readonly REST but here it was a quick way to show
// that OpenAPI spec changed and that get/list still work.

// rest implements a RESTStorage for API services against etcd
type REST struct {
	store *genericregistry.Store
}

var _ rest.Getter = &REST{}
var _ rest.Lister = &REST{}
var _ rest.Storage = &REST{}
var _ rest.Scoper = &REST{}
var _ rest.SingularNameProvider = &REST{}

func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*REST, error) {
	strategy := NewStrategy(scheme)

	store := &genericregistry.Store{
		NewFunc:                   func() runtime.Object { return &googlesheets.Datasource{} },
		NewListFunc:               func() runtime.Object { return &googlesheets.DatasourceList{} },
		PredicateFunc:             MatchDatasource,
		DefaultQualifiedResource:  googlesheets.Resource("datasources"),
		SingularQualifiedResource: googlesheets.Resource("datasource"),

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,

		// TODO: define table converter that exposes more than name/creation timestamp
		TableConvertor: rest.NewDefaultTableConvertor(googlesheets.Resource("datasources")),
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store: store}, nil
}

// Storage

func (dsREST REST) New() runtime.Object {
	return &googlesheets.Datasource{}
}

func (dsREST REST) Destroy() {
}

// Getter

func (dsREST REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return dsREST.store.Get(ctx, name, options)
}

func (dsREST REST) NamespaceScoped() bool {
	return true
}

func (dsREST REST) GetSingularName() string {
	return "datasource"
}

// List

func (dsREST REST) NewList() runtime.Object {
	return &googlesheets.DatasourceList{}
}

func (dsREST REST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	return dsREST.store.List(ctx, options)
}

func (dsREST REST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return dsREST.store.ConvertToTable(ctx, object, tableOptions)
}
