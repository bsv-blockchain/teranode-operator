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

// AssetSpec defines the desired state of Asset
type AssetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	NodeSelector       map[string]string            `json:"nodeSelector,omitempty"`
	Tolerations        *[]corev1.Toleration         `json:"tolerations,omitempty"`
	Affinity           *corev1.Affinity             `json:"affinity,omitempty"`
	Resources          *corev1.ResourceRequirements `json:"resources,omitempty"`
	GrpcIngress        *IngressDef                  `json:"grpcIngress,omitempty"`
	HTTPIngress        *IngressDef                  `json:"httpIngress,omitempty"`
	HTTPSIngress       *IngressDef                  `json:"httpsIngress,omitempty"`
	Image              string                       `json:"image,omitempty"`
	ImagePullPolicy    corev1.PullPolicy            `json:"imagePullPolicy,omitempty"`
	ServiceAccount     string                       `json:"serviceAccount,omitempty"`
	ConfigMapName      string                       `json:"configMapName,omitempty"`
	ServiceAnnotations map[string]string            `json:"serviceAnnotations,omitempty"`
	Replicas           *int32                       `json:"replicas,omitempty"`
	Command            []string                     `json:"command,omitempty"`
	Args               []string                     `json:"args,omitempty"`
}

// AssetStatus defines the observed state of Asset
type AssetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Asset is the Schema for the assets API
type Asset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AssetSpec   `json:"spec,omitempty"`
	Status AssetStatus `json:"status,omitempty"`
}

func (in *Asset) NodeSelector() map[string]string {
	return in.Spec.NodeSelector
}

func (in *Asset) Tolerations() *[]corev1.Toleration {
	return in.Spec.Tolerations
}

func (in *Asset) Affinity() *corev1.Affinity {
	return in.Spec.Affinity
}

func (in *Asset) Resources() *corev1.ResourceRequirements {
	return in.Spec.Resources
}

func (in *Asset) Image() string {
	return in.Spec.Image
}

func (in *Asset) ImagePullPolicy() corev1.PullPolicy {
	return in.Spec.ImagePullPolicy
}

func (in *Asset) ServiceAccountName() string {
	return in.Spec.ServiceAccount
}

func (in *Asset) Replicas() *int32 {
	return in.Spec.Replicas
}

func (in *Asset) ConfigMapName() string {
	return in.Spec.ConfigMapName
}

func (in *Asset) Command() []string {
	return in.Spec.Command
}

func (in *Asset) Args() []string {
	return in.Spec.Args
}

//+kubebuilder:object:root=true

// AssetList contains a list of Asset
type AssetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Asset `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Asset{}, &AssetList{})
}
