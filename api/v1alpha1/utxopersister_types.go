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

// UtxoPersisterSpec defines the desired state of UtxoPersister
type UtxoPersisterSpec struct {
	NodeSelector    map[string]string            `json:"nodeSelector,omitempty"`
	Tolerations     *[]corev1.Toleration         `json:"tolerations,omitempty"`
	Affinity        *corev1.Affinity             `json:"affinity,omitempty"`
	Resources       *corev1.ResourceRequirements `json:"resources,omitempty"`
	Image           string                       `json:"image,omitempty"`
	ImagePullPolicy corev1.PullPolicy            `json:"imagePullPolicy,omitempty"`
	ServiceAccount  string                       `json:"serviceAccount,omitempty"`
	ConfigMapName   string                       `json:"configMapName,omitempty"`
	Replicas        *int32                       `json:"replicas,omitempty"`
	Command         []string                     `json:"command,omitempty"`
	Args            []string                     `json:"args,omitempty"`
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

func (in *UtxoPersister) NodeSelector() map[string]string {
	return in.Spec.NodeSelector
}

func (in *UtxoPersister) Tolerations() *[]corev1.Toleration {
	return in.Spec.Tolerations
}

func (in *UtxoPersister) Affinity() *corev1.Affinity {
	return in.Spec.Affinity
}

func (in *UtxoPersister) Resources() *corev1.ResourceRequirements {
	return in.Spec.Resources
}

func (in *UtxoPersister) Image() string {
	return in.Spec.Image
}

func (in *UtxoPersister) ImagePullPolicy() corev1.PullPolicy {
	return in.Spec.ImagePullPolicy
}

func (in *UtxoPersister) ServiceAccountName() string {
	return in.Spec.ServiceAccount
}

func (in *UtxoPersister) Replicas() *int32 {
	return in.Spec.Replicas
}

func (in *UtxoPersister) ConfigMapName() string {
	return in.Spec.ConfigMapName
}

func (in *UtxoPersister) Command() []string {
	return in.Spec.Command
}

func (in *UtxoPersister) Args() []string {
	return in.Spec.Args
}

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
