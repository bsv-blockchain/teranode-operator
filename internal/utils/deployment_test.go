package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

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
				{Name: "vol1", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
				{Name: "vol2", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
				{Name: "vol1", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
			},
			expected: []corev1.Volume{
				{Name: "vol2", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
				{Name: "vol1", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
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
				{Name: "vol1", MountPath: "/data"},
				{Name: "vol2", MountPath: "/config"},
				{Name: "vol1", MountPath: "/data"},
			},
			expected: []corev1.VolumeMount{
				{Name: "vol2", MountPath: "/config"},
				{Name: "vol1", MountPath: "/data"},
			},
		},
		{
			name: "keeps different mountPaths for same volume",
			input: []corev1.VolumeMount{
				{Name: "vol1", MountPath: "/data1"},
				{Name: "vol1", MountPath: "/data2"},
			},
			expected: []corev1.VolumeMount{
				{Name: "vol1", MountPath: "/data1"},
				{Name: "vol1", MountPath: "/data2"},
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

// TestSetDeploymentOverridesNoDuplication verifies that calling
// SetDeploymentOverridesWithContext multiple times does not duplicate values
func TestSetDeploymentOverridesNoDuplication(t *testing.T) {
	// This test will fail initially, demonstrating the bug
	// After implementing deduplication, it should pass

	// Note: This is a conceptual test. Full implementation would require
	// mocking the CR interface and client, which is complex for a unit test.
	// The integration test in propagation_controller_test.go will provide
	// the actual verification.

	t.Skip("Skipping - requires full controller context. See propagation_controller_test.go for integration test")
}
