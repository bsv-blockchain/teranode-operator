package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
)

const testVolumeName = "vol1"

func TestDeduplicateEnvVars(t *testing.T) {
	tests := []struct {
		name     string
		input    []corev1.EnvVar
		expected []corev1.EnvVar
	}{
		{
			name: "removes exact duplicates",
			input: []corev1.EnvVar{
				{Name: "VAR1", Value: "value1"},
				{Name: "VAR2", Value: "value2"},
				{Name: "VAR1", Value: "value1"},
			},
			expected: []corev1.EnvVar{
				{Name: "VAR2", Value: "value2"},
				{Name: "VAR1", Value: "value1"},
			},
		},
		{
			name: "keeps last occurrence for same name",
			input: []corev1.EnvVar{
				{Name: "VAR1", Value: "old"},
				{Name: "VAR1", Value: "new"},
			},
			expected: []corev1.EnvVar{
				{Name: "VAR1", Value: "new"},
			},
		},
		{
			name:     "handles empty slice",
			input:    []corev1.EnvVar{},
			expected: []corev1.EnvVar{},
		},
		{
			name: "preserves single element",
			input: []corev1.EnvVar{
				{Name: "VAR1", Value: "value1"},
			},
			expected: []corev1.EnvVar{
				{Name: "VAR1", Value: "value1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deduplicateEnvVars(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDeduplicateEnvFrom(t *testing.T) {
	tests := []struct {
		name     string
		input    []corev1.EnvFromSource
		expected []corev1.EnvFromSource
	}{
		{
			name: "removes duplicate ConfigMapRef",
			input: []corev1.EnvFromSource{
				{ConfigMapRef: &corev1.ConfigMapEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: "cm1"}}},
				{ConfigMapRef: &corev1.ConfigMapEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: "cm2"}}},
				{ConfigMapRef: &corev1.ConfigMapEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: "cm1"}}},
			},
			expected: []corev1.EnvFromSource{
				{ConfigMapRef: &corev1.ConfigMapEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: "cm2"}}},
				{ConfigMapRef: &corev1.ConfigMapEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: "cm1"}}},
			},
		},
		{
			name: "removes duplicate SecretRef",
			input: []corev1.EnvFromSource{
				{SecretRef: &corev1.SecretEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: "secret1"}}},
				{SecretRef: &corev1.SecretEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: "secret1"}}},
			},
			expected: []corev1.EnvFromSource{
				{SecretRef: &corev1.SecretEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: "secret1"}}},
			},
		},
		{
			name:     "handles empty slice",
			input:    []corev1.EnvFromSource{},
			expected: []corev1.EnvFromSource{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deduplicateEnvFrom(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDeduplicateImagePullSecrets(t *testing.T) {
	tests := []struct {
		name     string
		input    []corev1.LocalObjectReference
		expected []corev1.LocalObjectReference
	}{
		{
			name: "removes duplicates",
			input: []corev1.LocalObjectReference{
				{Name: "secret1"},
				{Name: "secret2"},
				{Name: "secret1"},
			},
			expected: []corev1.LocalObjectReference{
				{Name: "secret2"},
				{Name: "secret1"},
			},
		},
		{
			name:     "handles empty slice",
			input:    []corev1.LocalObjectReference{},
			expected: []corev1.LocalObjectReference{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deduplicateImagePullSecrets(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDeduplicateVolumes(t *testing.T) {
	tests := []struct {
		name     string
		input    []corev1.Volume
		expected []corev1.Volume
	}{
		{
			name: "removes duplicates by name",
			input: []corev1.Volume{
				{Name: testVolumeName, VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
				{Name: "vol2", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
				{Name: testVolumeName, VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
			},
			expected: []corev1.Volume{
				{Name: "vol2", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
				{Name: testVolumeName, VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
			},
		},
		{
			name:     "handles empty slice",
			input:    []corev1.Volume{},
			expected: []corev1.Volume{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deduplicateVolumes(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDeduplicateVolumeMounts(t *testing.T) {
	tests := []struct {
		name     string
		input    []corev1.VolumeMount
		expected []corev1.VolumeMount
	}{
		{
			name: "removes duplicates by name and mountPath",
			input: []corev1.VolumeMount{
				{Name: testVolumeName, MountPath: "/data"},
				{Name: "vol2", MountPath: "/config"},
				{Name: testVolumeName, MountPath: "/data"},
			},
			expected: []corev1.VolumeMount{
				{Name: "vol2", MountPath: "/config"},
				{Name: testVolumeName, MountPath: "/data"},
			},
		},
		{
			name: "keeps different mountPaths for same volume",
			input: []corev1.VolumeMount{
				{Name: testVolumeName, MountPath: "/data1"},
				{Name: testVolumeName, MountPath: "/data2"},
			},
			expected: []corev1.VolumeMount{
				{Name: testVolumeName, MountPath: "/data1"},
				{Name: testVolumeName, MountPath: "/data2"},
			},
		},
		{
			name:     "handles empty slice",
			input:    []corev1.VolumeMount{},
			expected: []corev1.VolumeMount{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deduplicateVolumeMounts(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// fakeService is a minimal TeranodeService implementation for unit-testing the
// deployment override helpers without spinning up envtest.
type fakeService struct {
	overrides *v1alpha1.DeploymentOverrides
}

func (f fakeService) DeploymentOverrides() *v1alpha1.DeploymentOverrides { return f.overrides }
func (f fakeService) Metadata() metav1.ObjectMeta                        { return metav1.ObjectMeta{} }

func newDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{Name: "main"}},
				},
			},
		},
	}
}

// SetDeploymentOverrides is the single helper every service funnels through, so
// verifying PriorityClassName here proves the CRD->Deployment step for all services.
func TestSetDeploymentOverridesSetsPriorityClassName(t *testing.T) {
	dep := newDeployment()
	svc := fakeService{overrides: &v1alpha1.DeploymentOverrides{PriorityClassName: "high-priority"}}

	SetDeploymentOverrides(nil, dep, svc)

	assert.Equal(t, "high-priority", dep.Spec.Template.Spec.PriorityClassName)
}

// An unset PriorityClassName must not clobber whatever the default spec already has.
func TestSetDeploymentOverridesEmptyPriorityClassNameDoesNotClobber(t *testing.T) {
	dep := newDeployment()
	dep.Spec.Template.Spec.PriorityClassName = "existing"
	svc := fakeService{overrides: &v1alpha1.DeploymentOverrides{}}

	SetDeploymentOverrides(nil, dep, svc)

	assert.Equal(t, "existing", dep.Spec.Template.Spec.PriorityClassName)
}

// Nil overrides is a no-op.
func TestSetDeploymentOverridesNilOverrides(t *testing.T) {
	dep := newDeployment()
	svc := fakeService{overrides: nil}

	SetDeploymentOverrides(nil, dep, svc)

	assert.Empty(t, dep.Spec.Template.Spec.PriorityClassName)
}
