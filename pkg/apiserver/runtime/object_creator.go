// SPDX-License-Identifier: AGPL-3.0-only

package runtime

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ runtime.ObjectCreater = (*objectCreator)(nil)

type objectCreator struct{}

func NewObjectCreator() *objectCreator {
	return &objectCreator{}
}

func (u objectCreator) New(gvr schema.GroupVersionKind) (runtime.Object, error) {
	uObj := &unstructured.Unstructured{}
	uObj.SetGroupVersionKind(gvr)
	return uObj, nil
}
