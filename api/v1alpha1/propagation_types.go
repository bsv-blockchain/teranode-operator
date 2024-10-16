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

// PropagationSpec defines the desired state of Propagation
type PropagationSpec struct {
	DeploymentOverrides *DeploymentOverrides `json:"deploymentOverrides,omitempty"`
	ServiceAnnotations  map[string]string    `json:"serviceAnnotations,omitempty"`
	DelveIngress        *IngressDef          `json:"delveIngress,omitempty"`
	QuicIngress         *IngressDef          `json:"quicIngress,omitempty"`
	GrpcIngress         *IngressDef          `json:"grpcIngress,omitempty"`
	HTTPIngress         *IngressDef          `json:"httpIngress,omitempty"`
	ProfilerIngress     *IngressDef          `json:"httpsIngress,omitempty"`
}

// PropagationStatus defines the observed state of Propagation
type PropagationStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Propagation is the Schema for the propagations API
type Propagation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PropagationSpec   `json:"spec,omitempty"`
	Status PropagationStatus `json:"status,omitempty"`
}

func (p *Propagation) DeploymentOverrides() *DeploymentOverrides {
	return p.Spec.DeploymentOverrides
}
func (p *Propagation) Metadata() metav1.ObjectMeta { return p.ObjectMeta }

//+kubebuilder:object:root=true

// PropagationList contains a list of Propagation
type PropagationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Propagation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Propagation{}, &PropagationList{})
}
