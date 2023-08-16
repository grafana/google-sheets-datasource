package registry

import (
	"github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
)

// rest implements a RESTStorage for API services against etcd
type REST struct {
	*genericregistry.Store
}

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
	return &REST{Store: store}, nil
}
