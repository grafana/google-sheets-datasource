// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/kubernetes/apiextensions-apiserver/blob/master/pkg/apis/apiextensions/install/install.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: The Kubernetes Authors.

package install

import (
	"github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets"
	googlesheetsV1 "github.com/grafana/google-sheets-datasource/pkg/apis/googlesheets/v1"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(googlesheets.AddToScheme(scheme))
	utilruntime.Must(googlesheetsV1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(googlesheetsV1.SchemeGroupVersion))
}
