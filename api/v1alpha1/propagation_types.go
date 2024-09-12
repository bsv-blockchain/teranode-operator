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

// PropagationSpec defines the desired state of Propagation
type PropagationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	NodeSelector       map[string]string            `json:"nodeSelector,omitempty"`
	Tolerations        *[]corev1.Toleration         `json:"tolerations,omitempty"`
	Affinity           *corev1.Affinity             `json:"affinity,omitempty"`
	Resources          *corev1.ResourceRequirements `json:"resources,omitempty"`
	DelveIngress       *IngressDef                  `json:"delveIngress,omitempty"`
	QuicIngress        *IngressDef                  `json:"quicIngress,omitempty"`
	GrpcIngress        *IngressDef                  `json:"grpcIngress,omitempty"`
	HttpIngress        *IngressDef                  `json:"httpIngress,omitempty"`
	ProfilerIngress    *IngressDef                  `json:"httpsIngress,omitempty"`
	Image              string                       `json:"image,omitempty"`
	ImagePullPolicy    corev1.PullPolicy            `json:"imagePullPolicy,omitempty"`
	ServiceAccount     string                       `json:"serviceAccount,omitempty"`
	ConfigMapName      string                       `json:"configMapName,omitempty"`
	ServiceAnnotations map[string]string            `json:"serviceAnnotations,omitempty"`
	Replicas           *int32                       `json:"replicas,omitempty"`
	Command            []string                     `json:"command,omitempty"`
	Args               []string                     `json:"args,omitempty"`
}

// PropagationStatus defines the observed state of Propagation
type PropagationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Propagation is the Schema for the propagations API
type Propagation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PropagationSpec   `json:"spec,omitempty"`
	Status PropagationStatus `json:"status,omitempty"`
}

func (in *Propagation) NodeSelector() map[string]string {
	return in.Spec.NodeSelector
}

func (in *Propagation) Tolerations() *[]corev1.Toleration {
	return in.Spec.Tolerations
}

func (in *Propagation) Affinity() *corev1.Affinity {
	return in.Spec.Affinity
}

func (in *Propagation) Resources() *corev1.ResourceRequirements {
	return in.Spec.Resources
}

func (in *Propagation) Image() string {
	return in.Spec.Image
}

func (in *Propagation) ImagePullPolicy() corev1.PullPolicy {
	return in.Spec.ImagePullPolicy
}

func (in *Propagation) ServiceAccountName() string {
	return in.Spec.ServiceAccount
}

func (in *Propagation) Replicas() *int32 {
	return in.Spec.Replicas
}

func (in *Propagation) ConfigMapName() string {
	return in.Spec.ConfigMapName
}

func (in *Propagation) Command() []string {
	return in.Spec.Command
}

func (in *Propagation) Args() []string {
	return in.Spec.Args
}

//+kubebuilder:object:root=true

// PropagationList contains a list of Propagation
type PropagationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Propagation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Propagation{}, &PropagationList{})
}
