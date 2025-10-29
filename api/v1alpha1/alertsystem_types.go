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

// AlertSystemSpec defines the desired state of AlertSystem
type AlertSystemSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DeploymentOverrides *DeploymentOverrides `json:"deploymentOverrides,omitempty"`
}

// AlertSystemStatus defines the observed state of AlertSystem
type AlertSystemStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func (m *AlertSystem) DeploymentOverrides() *DeploymentOverrides {
	return m.Spec.DeploymentOverrides
}
func (m *AlertSystem) Metadata() metav1.ObjectMeta { return m.ObjectMeta }

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AlertSystem is the Schema for the alertsystems API
type AlertSystem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlertSystemSpec   `json:"spec,omitempty"`
	Status AlertSystemStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AlertSystemList contains a list of AlertSystem
type AlertSystemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []AlertSystem `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AlertSystem{}, &AlertSystemList{})
}
