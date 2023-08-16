package storage

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"

	pluginRuntime "github.com/grafana/google-sheets-datasource/pkg/apiserver/runtime"
)

type PluginResourceStorage struct {
	Query *SubresourceStreamerREST
}

func NewStorage(resource schema.GroupResource, singularResource schema.GroupResource, kind, listKind schema.GroupVersionKind, optsGetter generic.RESTOptionsGetter, tableConvertor rest.TableConvertor, typer pluginRuntime.ObjectTyper) *PluginResourceStorage {
	var storage PluginResourceStorage
	return &storage
}
