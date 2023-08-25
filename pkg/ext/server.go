package ext

import (
	"context"
	"io"
	"net/http"

	"github.com/grafana/kindsys"
	"k8s.io/kube-openapi/pkg/spec3"
)

// Alternative "Grafana Resource Definition" hook
// An api server will be created with support for the kind and the various
type APIServiceHooks struct {
	// This defines the apiVersion + Kind
	Kind kindsys.ResourceKind

	// Called before creating a new resource (admission/mutation controller)
	// this can error or mutate
	BeforeAdd func(ctx context.Context, obj kindsys.Resource) (kindsys.Resource, error)

	// Called before updating a resource  (admission controller)
	// this can error or mutate
	BeforeUpdate func(ctx context.Context, newObj kindsys.Resource, oldObj kindsys.Resource) (kindsys.Resource, error)

	// Called before deleting a resource
	// this can error
	// ??? is this necessary -- finalizers seem like the real thing
	// BeforeDelete func(ctx context.Context, obj Resource) error

	// This is called when initialized -- the endpoints will be added to the api server
	// the OpenAPI specs will be exposed in the public API
	PluginRouteHandlers []PluginRouteHandler
}

// This is used to implement dynamic sub-resources like pods/x/logs
type PluginRouteHandler struct {
	Level   RawAPILevel      // group+version | namespace | resource
	Slug    string           // added to the appropriate level
	Spec    *spec3.PathProps // Exposed in the open api service discovery
	Handler http.HandlerFunc // when Level = resource, the resource will be available in context
}

// This allows access to resources for API handlers
type ResourceGetter = func(ctx context.Context, ns string, name string) (kindsys.Resource, error)

// This is used to answer raw API requests like /logs
type StreamingResponse = func(ctx context.Context, apiVersion, acceptHeader string) (
	stream io.ReadCloser, flush bool, mimeType string, err error)

type RawAPILevel int8

const (
	RawAPILevelResource RawAPILevel = iota
	RawAPILevelNamespace
	RawAPILevelGroupVersion
)
