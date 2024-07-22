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

// BlockPersisterSpec defines the desired state of BlockPersister
type BlockPersisterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	NodeSelector     map[string]string            `json:"nodeSelector,omitempty"`
	Tolerations      *[]corev1.Toleration         `json:"tolerations,omitempty"`
	Affinity         *corev1.Affinity             `json:"affinity,omitempty"`
	Resources        *corev1.ResourceRequirements `json:"resources,omitempty"`
	StorageClass     string                       `json:"storageClass,omitempty"`
	StorageResources *corev1.ResourceRequirements `json:"storageResources,omitempty"`
	Image            string                       `json:"image,omitempty"`
	ImagePullPolicy  corev1.PullPolicy            `json:"imagePullPolicy,omitempty"`
	ServiceAccount   string                       `json:"serviceAccount,omitempty"`
	ConfigMapName    string                       `json:"configMapName,omitempty"`
}

// BlockPersisterStatus defines the observed state of BlockPersister
type BlockPersisterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BlockPersister is the Schema for the blockpersisters API
type BlockPersister struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BlockPersisterSpec   `json:"spec,omitempty"`
	Status BlockPersisterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BlockPersisterList contains a list of BlockPersister
type BlockPersisterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BlockPersister `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BlockPersister{}, &BlockPersisterList{})
}
