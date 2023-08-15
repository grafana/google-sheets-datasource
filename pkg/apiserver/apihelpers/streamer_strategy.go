// SPDX-License-Identifier: AGPL-3.0-only

package apihelpers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage/names"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

type StreamerStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

var _ rest.RESTCreateStrategy = (*StreamerStrategy)(nil)
var _ rest.RESTUpdateStrategy = (*StreamerStrategy)(nil)
var _ rest.RESTDeleteStrategy = (*StreamerStrategy)(nil)
var _ rest.ResetFieldsStrategy = (*StreamerStrategy)(nil)

func NewStreamerStrategy(typer runtime.ObjectTyper) StreamerStrategy {
	return StreamerStrategy{
		ObjectTyper:   typer,
		NameGenerator: names.SimpleNameGenerator,
	}
}

func (StreamerStrategy) NamespaceScoped() bool {
	return true
}

func (StreamerStrategy) Canonicalize(obj runtime.Object) {}

func (StreamerStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {}

func (StreamerStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (StreamerStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (StreamerStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {}

func (StreamerStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

func (StreamerStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (StreamerStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (StreamerStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (StreamerStrategy) PrepareForDelete(ctx context.Context, obj runtime.Object) {}

func (StreamerStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return map[fieldpath.APIVersion]*fieldpath.Set{}
}
