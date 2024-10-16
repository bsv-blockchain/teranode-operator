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

// LegacySpec defines the desired state of Legacy
type LegacySpec struct {
	DeploymentOverrides *DeploymentOverrides `json:"deploymentOverrides,omitempty"`
}

// LegacyStatus defines the observed state of Legacy
type LegacyStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Legacy is the Schema for the legacies API
type Legacy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LegacySpec   `json:"spec,omitempty"`
	Status LegacyStatus `json:"status,omitempty"`
}

func (l *Legacy) DeploymentOverrides() *DeploymentOverrides {
	return l.Spec.DeploymentOverrides
}
func (l *Legacy) Metadata() metav1.ObjectMeta { return l.ObjectMeta }

//+kubebuilder:object:root=true

// LegacyList contains a list of Legacy
type LegacyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Legacy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Legacy{}, &LegacyList{})
}
