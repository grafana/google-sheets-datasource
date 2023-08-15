// SPDX-License-Identifier: AGPL-3.0-only

package runtime

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type objectConvertor struct{}

var _ runtime.ObjectConvertor = (*objectConvertor)(nil)

func NewObjectConvertor() *objectConvertor {
	return &objectConvertor{}
}

// TODO: implement this
func (c *objectConvertor) Convert(in, out, context interface{}) error {
	out = in
	return nil
}

func (c *objectConvertor) ConvertToVersion(in runtime.Object, gv runtime.GroupVersioner) (out runtime.Object, err error) {
	return in.DeepCopyObject(), nil
}

func (c *objectConvertor) ConvertFieldLabel(gvk schema.GroupVersionKind, label, value string) (string, string, error) {
	return label, value, nil
}
