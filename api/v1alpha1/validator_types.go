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

// ValidatorSpec defines the desired state of Validator
type ValidatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

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

// ValidatorStatus defines the observed state of Validator
type ValidatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Validator is the Schema for the validators API
type Validator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ValidatorSpec   `json:"spec,omitempty"`
	Status ValidatorStatus `json:"status,omitempty"`
}

func (in *Validator) NodeSelector() map[string]string {
	return in.Spec.NodeSelector
}

func (in *Validator) Tolerations() *[]corev1.Toleration {
	return in.Spec.Tolerations
}

func (in *Validator) Affinity() *corev1.Affinity {
	return in.Spec.Affinity
}

func (in *Validator) Resources() *corev1.ResourceRequirements {
	return in.Spec.Resources
}

func (in *Validator) Image() string {
	return in.Spec.Image
}

func (in *Validator) ImagePullPolicy() corev1.PullPolicy {
	return in.Spec.ImagePullPolicy
}

func (in *Validator) ServiceAccountName() string {
	return in.Spec.ServiceAccount
}

func (in *Validator) Replicas() *int32 {
	return in.Spec.Replicas
}

func (in *Validator) ConfigMapName() string {
	return in.Spec.ConfigMapName
}

func (in *Validator) Command() []string {
	return in.Spec.Command
}

func (in *Validator) Args() []string {
	return in.Spec.Args
}

//+kubebuilder:object:root=true

// ValidatorList contains a list of Validator
type ValidatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Validator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Validator{}, &ValidatorList{})
}
