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

// UtxoPersisterSpec defines the desired state of UtxoPersister
type UtxoPersisterSpec struct {
	DeploymentOverrides *DeploymentOverrides `json:"deploymentOverrides,omitempty"`
}

// UtxoPersisterStatus defines the observed state of UtxoPersister
type UtxoPersisterStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// UtxoPersister is the Schema for the utxopersisters API
type UtxoPersister struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UtxoPersisterSpec   `json:"spec,omitempty"`
	Status UtxoPersisterStatus `json:"status,omitempty"`
}

func (up *UtxoPersister) DeploymentOverrides() *DeploymentOverrides {
	return up.Spec.DeploymentOverrides
}
func (up *UtxoPersister) Metadata() metav1.ObjectMeta { return up.ObjectMeta }

//+kubebuilder:object:root=true

// UtxoPersisterList contains a list of UtxoPersister
type UtxoPersisterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UtxoPersister `json:"items"`
}

func init() {
	SchemeBuilder.Register(&UtxoPersister{}, &UtxoPersisterList{})
}
