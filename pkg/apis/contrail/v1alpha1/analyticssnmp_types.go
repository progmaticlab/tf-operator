package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AnalyticsSnmp is the Schema for the Analytics SNMP API.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=analyticssnmp,scope=Namespaced
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.status.active`
type AnalyticsSnmp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AnalyticsSnmpSpec   `json:"spec,omitempty"`
	Status AnalyticsSnmpStatus `json:"status,omitempty"`
}

// AnalyticsSnmpList contains a list of AnalyticsSnmp.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AnalyticsSnmpList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []AnalyticsSnmp `json:"items"`
}

// AnalyticsSnmpSpec is the Spec for the Analytics SNMP API.
// +k8s:openapi-gen=true
type AnalyticsSnmpSpec struct {
	CommonConfiguration  PodConfiguration           `json:"commonConfiguration,omitempty"`
	ServiceConfiguration AnalyticsSnmpConfiguration `json:"serviceConfiguration"`
}

// AnalyticsSnmpConfiguration is the Spec for the Analytics SNMP API.
// +k8s:openapi-gen=true
type AnalyticsSnmpConfiguration struct {
	Containers []*Container `json:"containers,omitempty"`
}

// AnalyticsSnmpStatus is the Status for the Analytics SNMP API.
// +k8s:openapi-gen=true
type AnalyticsSnmpStatus struct {
	Active *bool `json:"active,omitempty"`
}

func init() {
	SchemeBuilder.Register(&AnalyticsSnmp{}, &AnalyticsSnmpList{})
}
