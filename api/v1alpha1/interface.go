package v1alpha1

import corev1 "k8s.io/api/core/v1"

// +genclient
// +k8s:deepcopy-gen=false
// TeranodeService interface defines the methods used for all common configuration overrides of each API service
type TeranodeService interface {
	NodeSelector() map[string]string
	Tolerations() *[]corev1.Toleration
	Affinity() *corev1.Affinity
	Resources() *corev1.ResourceRequirements
	Image() string
	ImagePullPolicy() corev1.PullPolicy
	ServiceAccountName() string
	Replicas() *int32
	ConfigMapName() string
	Command() []string
	Args() []string
}
