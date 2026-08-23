package controller

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"

	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
)

// mergeDeploymentOverrides is generic (reflection over every field), so any set field
// on the cluster overrides propagates to the child. This guards that behavior.
func TestMergeDeploymentOverridesPropagatesSetFields(t *testing.T) {
	target := &teranodev1alpha1.DeploymentOverrides{}
	taints := []corev1.Taint{{Key: "dedicated", Value: "teranode", Effect: corev1.TaintEffectNoSchedule}}
	source := &teranodev1alpha1.DeploymentOverrides{
		PriorityClassName: "high-priority",
		ServiceAccount:    "svc-acct",
		Replicas:          ptr.To(int32(4)),
		Taints:            &taints, // Taints was silently dropped by the old field-by-field merge
	}

	mergeDeploymentOverrides(target, source)

	if target.PriorityClassName != "high-priority" {
		t.Errorf("PriorityClassName not propagated, got %q", target.PriorityClassName)
	}
	if target.ServiceAccount != "svc-acct" {
		t.Errorf("ServiceAccount not propagated, got %q", target.ServiceAccount)
	}
	if target.Replicas == nil || *target.Replicas != 4 {
		t.Errorf("Replicas not propagated, got %v", target.Replicas)
	}
	if target.Taints == nil {
		t.Error("Taints not propagated (regression: old merge omitted this field)")
	}
}

// Fields left unset on the source must not overwrite existing values on the target.
func TestMergeDeploymentOverridesDoesNotClobberUnsetFields(t *testing.T) {
	existingReplicas := ptr.To(int32(3))
	target := &teranodev1alpha1.DeploymentOverrides{
		PriorityClassName: "keep-me",
		Replicas:          existingReplicas,
	}
	source := &teranodev1alpha1.DeploymentOverrides{
		Image: "new-image", // only Image set
	}

	mergeDeploymentOverrides(target, source)

	if target.PriorityClassName != "keep-me" {
		t.Errorf("unset source field clobbered target PriorityClassName, got %q", target.PriorityClassName)
	}
	if target.Replicas != existingReplicas {
		t.Errorf("unset source field clobbered target Replicas, got %v", target.Replicas)
	}
	if target.Image != "new-image" {
		t.Errorf("set source field not applied, got %q", target.Image)
	}
}

// A nil source is a safe no-op.
func TestMergeDeploymentOverridesNilSource(t *testing.T) {
	target := &teranodev1alpha1.DeploymentOverrides{PriorityClassName: "keep-me"}

	mergeDeploymentOverrides(target, nil)

	if target.PriorityClassName != "keep-me" {
		t.Errorf("nil source mutated target, got %q", target.PriorityClassName)
	}
}
