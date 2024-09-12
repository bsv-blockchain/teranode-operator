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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BlockchainSpec defines the desired state of Blockchain
type BlockchainSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Resources       *corev1.ResourceRequirements `json:"resources,omitempty"`
	Image           string                       `json:"image,omitempty"`
	ImagePullPolicy corev1.PullPolicy            `json:"imagePullPolicy,omitempty"`
	NodeSelector    map[string]string            `json:"nodeSelector,omitempty"`
	Tolerations     *[]corev1.Toleration         `json:"tolerations,omitempty"`
	Affinity        *corev1.Affinity             `json:"affinity,omitempty"`
	ServiceAccount  string                       `json:"serviceAccount,omitempty"`
	ConfigMapName   string                       `json:"configMapName,omitempty"`
	Replicas        *int32                       `json:"replicas,omitempty"`
	Command         []string                     `json:"command,omitempty"`
	Args            []string                     `json:"args,omitempty"`
}

// BlockchainStatus defines the observed state of Blockchain
type BlockchainStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions,omitempty"`
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

func (in *Blockchain) NodeSelector() map[string]string {
	return in.Spec.NodeSelector
}

func (in *Blockchain) Tolerations() *[]corev1.Toleration {
	return in.Spec.Tolerations
}

func (in *Blockchain) Affinity() *corev1.Affinity {
	return in.Spec.Affinity
}

func (in *Blockchain) Resources() *corev1.ResourceRequirements {
	return in.Spec.Resources
}

func (in *Blockchain) Image() string {
	return in.Spec.Image
}

func (in *Blockchain) ImagePullPolicy() corev1.PullPolicy {
	return in.Spec.ImagePullPolicy
}

func (in *Blockchain) ServiceAccountName() string {
	return in.Spec.ServiceAccount
}

func (in *Blockchain) Replicas() *int32 {
	return in.Spec.Replicas
}

func (in *Blockchain) ConfigMapName() string {
	return in.Spec.ConfigMapName
}

func (in *Blockchain) Command() []string {
	return in.Spec.Command
}

func (in *Blockchain) Args() []string {
	return in.Spec.Args
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
