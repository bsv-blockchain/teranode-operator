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

// BlockValidatorSpec defines the desired state of BlockValidator
type BlockValidatorSpec struct {
	DeploymentOverrides *DeploymentOverrides `json:"deploymentOverrides,omitempty"`
}

// BlockValidatorStatus defines the observed state of BlockValidator
type BlockValidatorStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BlockValidator is the Schema for the blockvalidators API
type BlockValidator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BlockValidatorSpec   `json:"spec,omitempty"`
	Status BlockValidatorStatus `json:"status,omitempty"`
}

func (bv *BlockValidator) DeploymentOverrides() *DeploymentOverrides {
	return bv.Spec.DeploymentOverrides
}
func (bv *BlockValidator) Metadata() metav1.ObjectMeta { return bv.ObjectMeta }

//+kubebuilder:object:root=true

// BlockValidatorList contains a list of BlockValidator
type BlockValidatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BlockValidator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BlockValidator{}, &BlockValidatorList{})
}
