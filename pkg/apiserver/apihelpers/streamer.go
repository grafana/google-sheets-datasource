package apihelpers

import (
	"context"
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"
)

type inputStreamFn = func(ctx context.Context, apiVersion, acceptHeader string) (stream io.ReadCloser, flush bool, contentType string, err error)

type SubresourceStreamer struct {
	fn inputStreamFn
}

var _ rest.ResourceStreamer = (*SubresourceStreamer)(nil)

func (s *SubresourceStreamer) SetInputStream(fn inputStreamFn) {
	s.fn = fn
}

func (s *SubresourceStreamer) GetObjectKind() schema.ObjectKind {
	return schema.EmptyObjectKind
}

func (s *SubresourceStreamer) DeepCopyObject() runtime.Object {
	panic("SubresourceStreamer does not implement DeepCopyObject")
}

func (s *SubresourceStreamer) InputStream(ctx context.Context, apiVersion, acceptHeader string) (stream io.ReadCloser, flush bool, contentType string, err error) {
	if s.fn == nil {
		return nil, false, "", fmt.Errorf("subresource streamer not initialized with input stream function")
	}
	return s.fn(ctx, apiVersion, acceptHeader)
}
