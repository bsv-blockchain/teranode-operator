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
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CoinbaseSpec defines the desired state of Coinbase
type CoinbaseSpec struct {
	GrpcIngress         *IngressDef          `json:"grpcIngress,omitempty"`
	DeploymentOverrides *DeploymentOverrides `json:"deploymentOverrides,omitempty"`
}

// CoinbaseStatus defines the observed state of Coinbase
type CoinbaseStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Coinbase is the Schema for the coinbases API
type Coinbase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CoinbaseSpec   `json:"spec,omitempty"`
	Status CoinbaseStatus `json:"status,omitempty"`
}

func (cb *Coinbase) DeploymentOverrides() *DeploymentOverrides {
	return cb.Spec.DeploymentOverrides
}
func (cb *Coinbase) Metadata() metav1.ObjectMeta { return cb.ObjectMeta }

//+kubebuilder:object:root=true

// CoinbaseList contains a list of Coinbase
type CoinbaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Coinbase `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Coinbase{}, &CoinbaseList{})
}
