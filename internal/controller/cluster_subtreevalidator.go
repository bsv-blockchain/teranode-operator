package controller

import (
	"fmt"

	"github.com/go-logr/logr"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
)

// ReconcileSubtreeValidator is the cluster subtreeValidator reconciler
func (r *ClusterReconciler) ReconcileSubtreeValidator(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	subtreeValidator := teranodev1alpha1.SubtreeValidator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-subtreevalidator", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels("subtree-validator"),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.SubtreeValidator.Enabled || (cluster.Spec.Enabled != nil && !*cluster.Spec.Enabled) {
		namespacedName := types.NamespacedName{
			Name:      subtreeValidator.Name,
			Namespace: subtreeValidator.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &subtreeValidator)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &subtreeValidator)
		return true, err
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &subtreeValidator, func() error {
		return r.updateSubtreeValidator(&subtreeValidator, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

//nolint:gocognit,gocyclo // Function complexity is inherent to handling multiple override cases
func (r *ClusterReconciler) updateSubtreeValidator(subtreeValidator *teranodev1alpha1.SubtreeValidator, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, subtreeValidator, r.Scheme)
	if err != nil {
		return err
	}

	// Only set defaults if this is a new resource (no spec configured yet)
	if subtreeValidator.Spec.DeploymentOverrides == nil && cluster.Spec.SubtreeValidator.Spec == nil {
		subtreeValidator.Spec = *defaultSubtreeValidatorSpec()
	}

	// Selectively merge cluster spec - only override fields that are explicitly set
	//nolint:nestif // Nested conditions required for selective field merging
	if cluster.Spec.SubtreeValidator.Spec != nil {
		clusterSpec := cluster.Spec.SubtreeValidator.Spec

		// Merge service-specific fields
		if clusterSpec.PodTemplateAnnotations != nil {
			subtreeValidator.Spec.PodTemplateAnnotations = clusterSpec.PodTemplateAnnotations
		}

		// Merge deployment overrides selectively
		if clusterSpec.DeploymentOverrides != nil {
			if subtreeValidator.Spec.DeploymentOverrides == nil {
				subtreeValidator.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
			}
			mergeDeploymentOverrides(subtreeValidator.Spec.DeploymentOverrides, clusterSpec.DeploymentOverrides)
		}
	}

	// Apply cluster-level defaults (only if not already set)
	if subtreeValidator.Spec.DeploymentOverrides == nil {
		subtreeValidator.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" {
		subtreeValidator.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}
	// Always apply cluster-level ImagePullSecrets (they override or are the default)
	if cluster.Spec.ImagePullSecrets != nil {
		subtreeValidator.Spec.DeploymentOverrides.ImagePullSecrets = cluster.Spec.ImagePullSecrets
	}

	return nil
}

func defaultSubtreeValidatorSpec() *teranodev1alpha1.SubtreeValidatorSpec {
	return &teranodev1alpha1.SubtreeValidatorSpec{}
}
