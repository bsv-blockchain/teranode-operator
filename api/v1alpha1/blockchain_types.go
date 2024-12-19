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

// BlockchainSpec defines the desired state of Blockchain
type BlockchainSpec struct {
	DeploymentOverrides *DeploymentOverrides        `json:"deploymentOverrides,omitempty"`
	FiniteStateMachine  *FiniteStateMachineSettings `json:"finiteStateMachine,omitempty"`
}

// FiniteStateMachineSettings defines the configuration of the FSM
type FiniteStateMachineSettings struct {
	Enabled      bool             `json:"enabled,omitempty"`
	Host         string           `json:"host,omitempty"`
	PollInterval *metav1.Duration `json:"pollInterval,omitempty"`
}

// BlockchainStatus defines the observed state of Blockchain
type BlockchainStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	FSMState   string             `json:"fsmState,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:rbac:groups="",resources=endpoints;configmaps;services;secrets;persistentvolumeclaims,verbs=get;create;update;list;watch
//+kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;update;create;list;watch
//+kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses,verbs=get;update;create;list;watch

// Blockchain is the Schema for the blockchains API
type Blockchain struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BlockchainSpec   `json:"spec,omitempty"`
	Status BlockchainStatus `json:"status,omitempty"`
}

func (bl *Blockchain) DeploymentOverrides() *DeploymentOverrides {
	return bl.Spec.DeploymentOverrides
}
func (bl *Blockchain) Metadata() metav1.ObjectMeta {
	return bl.ObjectMeta
}

//+kubebuilder:object:root=true

// BlockchainList contains a list of Blockchain
type BlockchainList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Blockchain `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Blockchain{}, &BlockchainList{})
}
