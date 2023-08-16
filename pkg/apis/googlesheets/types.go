package googlesheets

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Datasource
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Datasource struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DatasourceSpec `json:"spec,omitempty"`
	// +optional
	Status DatasourceStatus `json:"status,omitempty"`
}

// DatasourceList
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DatasourceList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Datasource `json:"items"`
}

// DatasourceSpec defines the desired state of Datasource
type DatasourceSpec struct {
	AuthType           string `json:"authType"` // jwt | key | gce
	APIKey             string `json:"apiKey"`
	DefaultProject     string `json:"defaultProject"`
	JWT                string `json:"jwt"`
	ClientEmail        string `json:"clientEmail"`
	TokenURI           string `json:"tokenUri"`
	AuthenticationType string `json:"authenticationType"`
	PrivateKeyPath     string `json:"privateKeyPath"`
	PrivateKey         string `json:"privateKey"` // `json:"-"`
}

// DatasourceStatus defines the observed state of Datasource
type DatasourceStatus struct {
}
