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

// AssetSpec defines the desired state of Asset
type AssetSpec struct {
	DeploymentOverrides *DeploymentOverrides `json:"deploymentOverrides,omitempty"`
	GrpcIngress         *IngressDef          `json:"grpcIngress,omitempty"`
	HTTPIngress         *IngressDef          `json:"httpIngress,omitempty"`
	HTTPSIngress        *IngressDef          `json:"httpsIngress,omitempty"`
}

// AssetStatus defines the observed state of Asset
type AssetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// Replicas is the number of actual replicas of the asset deployment
	Replicas int32 `json:"replicas,omitempty"`
	// Selector is the label selector for pods corresponding to this asset deployment
	Selector string `json:"selector,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:subresource:scale:specpath=.spec.deploymentOverrides.replicas,statuspath=.status.replicas,selectorpath=.status.selector

// Asset is the Schema for the assets API
type Asset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AssetSpec   `json:"spec,omitempty"`
	Status AssetStatus `json:"status,omitempty"`
}

func (a *Asset) DeploymentOverrides() *DeploymentOverrides {
	return a.Spec.DeploymentOverrides
}

func (a *Asset) Metadata() metav1.ObjectMeta {
	return a.ObjectMeta
}

//+kubebuilder:object:root=true

// AssetList contains a list of Asset
type AssetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Asset `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Asset{}, &AssetList{})
}
