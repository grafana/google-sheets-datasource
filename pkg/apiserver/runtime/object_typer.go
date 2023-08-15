// SPDX-License-Identifier: AGPL-3.0-only

package runtime

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ runtime.ObjectTyper = (*ObjectTyper)(nil)

type ObjectTyper struct{}

func NewObjectTyper() *ObjectTyper {
	return &ObjectTyper{}
}

func (dt ObjectTyper) ObjectKinds(obj runtime.Object) ([]schema.GroupVersionKind, bool, error) {
	return []schema.GroupVersionKind{obj.GetObjectKind().GroupVersionKind()}, false, nil
}

func (dt ObjectTyper) Recognizes(gvk schema.GroupVersionKind) bool {
	return true
}
