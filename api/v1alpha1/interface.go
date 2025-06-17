package v1alpha1

import (
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen=false

// TeranodeService defines the methods used for all common configuration overrides of each API service
type TeranodeService interface {
	DeploymentOverrides() *DeploymentOverrides
	Metadata() metav1.ObjectMeta
}

// DeploymentOverrides defines all the overrides for the deployment of each service
type DeploymentOverrides struct {
	NodeSelector       map[string]string            `json:"nodeSelector,omitempty"`
	Tolerations        *[]corev1.Toleration         `json:"tolerations,omitempty"`
	Taints             *[]corev1.Taint              `json:"taints,omitempty"`
	Affinity           *corev1.Affinity             `json:"affinity,omitempty"`
	PodAntiAffinity    *corev1.PodAntiAffinity      `json:"podAntiAffinity,omitempty"`
	Resources          *corev1.ResourceRequirements `json:"resources,omitempty"`
	Image              string                       `json:"image,omitempty"`
	ImagePullPolicy    corev1.PullPolicy            `json:"imagePullPolicy,omitempty"`
	ServiceAccount     string                       `json:"serviceAccount,omitempty"`
	ConfigMapName      string                       `json:"configMapName,omitempty"`
	ServiceAnnotations map[string]string            `json:"serviceAnnotations,omitempty"`
	Replicas           *int32                       `json:"replicas,omitempty"`
	Command            []string                     `json:"command,omitempty"`
	Args               []string                     `json:"args,omitempty"`
	Strategy           *v1.DeploymentStrategy       `json:"strategy,omitempty"`
	Env                []corev1.EnvVar              `json:"env,omitempty"`
	EnvFrom            []corev1.EnvFromSource       `json:"envFrom,omitempty"`
}
