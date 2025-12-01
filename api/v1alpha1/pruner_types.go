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

// PrunerSpec defines the desired state of Pruner
type PrunerSpec struct {
	DeploymentOverrides *DeploymentOverrides `json:"deploymentOverrides,omitempty"`
}

// PrunerStatus defines the observed state of Pruner
type PrunerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Pruner is the Schema for the pruners API
type Pruner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PrunerSpec   `json:"spec,omitempty"`
	Status PrunerStatus `json:"status,omitempty"`
}

func (p *Pruner) DeploymentOverrides() *DeploymentOverrides {
	return p.Spec.DeploymentOverrides
}

func (p *Pruner) Metadata() metav1.ObjectMeta {
	return p.ObjectMeta
}

//+kubebuilder:object:root=true

// PrunerList contains a list of Pruner
type PrunerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Pruner `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pruner{}, &PrunerList{})
}
