/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
//nolint:godox // Kubebuilder-generated scaffolding comment
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RPCSpec defines the desired state of RPC
type RPCSpec struct {
	DeploymentOverrides *DeploymentOverrides `json:"deploymentOverrides,omitempty"`
}

// RPCStatus defines the observed state of RPC
type RPCStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func (m *RPC) DeploymentOverrides() *DeploymentOverrides {
	return m.Spec.DeploymentOverrides
}
func (m *RPC) Metadata() metav1.ObjectMeta { return m.ObjectMeta }

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RPC is the Schema for the rpcs API
type RPC struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RPCSpec   `json:"spec,omitempty"`
	Status RPCStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RPCList contains a list of RPC
type RPCList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []RPC `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RPC{}, &RPCList{})
}
