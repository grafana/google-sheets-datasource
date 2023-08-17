package registry

import (
	"context"
	"fmt"

	"github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	apifield "k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
)

// NewStrategy creates and returns a datasourceStrategy instance
func NewStrategy(typer runtime.ObjectTyper) datasourceStrategy {
	return datasourceStrategy{typer, names.SimpleNameGenerator}
}

// GetAttrs returns labels.Set, fields.Set, and error in case the given runtime.Object is not a Datasource
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	apiserver, ok := obj.(*googlesheets.Datasource)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Datasource")
	}
	return labels.Set(apiserver.ObjectMeta.Labels), SelectableFields(apiserver), nil
}

// MatchDatasource is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func MatchDatasource(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *googlesheets.Datasource) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

type datasourceStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (datasourceStrategy) NamespaceScoped() bool {
	return true
}

func (datasourceStrategy) PrepareForCreate(_ context.Context, _ runtime.Object) {
}

func (datasourceStrategy) PrepareForUpdate(_ context.Context, _, _ runtime.Object) {
}

func (datasourceStrategy) Validate(_ context.Context, _ runtime.Object) apifield.ErrorList {
	return apifield.ErrorList{}
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (datasourceStrategy) WarningsOnCreate(_ context.Context, _ runtime.Object) []string {
	return nil
}

func (datasourceStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (datasourceStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (datasourceStrategy) Canonicalize(_ runtime.Object) {
}

func (datasourceStrategy) ValidateUpdate(_ context.Context, _, _ runtime.Object) apifield.ErrorList {
	return apifield.ErrorList{}
}

// WarningsOnUpdate returns warnings for the given update.
func (datasourceStrategy) WarningsOnUpdate(_ context.Context, _, _ runtime.Object) []string {
	return nil
}
