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

// ReconcilePropagation is the cluster propagation reconciler
func (r *ClusterReconciler) ReconcilePropagation(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	propagation := teranodev1alpha1.Propagation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-propagation", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels("propagation"),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.Propagation.Enabled || (cluster.Spec.Enabled != nil && !*cluster.Spec.Enabled) {
		namespacedName := types.NamespacedName{
			Name:      propagation.Name,
			Namespace: propagation.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &propagation)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &propagation)
		return true, err
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &propagation, func() error {
		return r.updatePropagation(&propagation, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

//nolint:gocognit,gocyclo // Function complexity is inherent to handling multiple override cases
func (r *ClusterReconciler) updatePropagation(propagation *teranodev1alpha1.Propagation, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, propagation, r.Scheme)
	if err != nil {
		return err
	}

	// Only set defaults if this is a new resource (no spec configured yet)
	if propagation.Spec.DeploymentOverrides == nil && cluster.Spec.Propagation.Spec == nil {
		propagation.Spec = *defaultPropagationSpec()
	}

	// Selectively merge cluster spec - only override fields that are explicitly set
	//nolint:nestif // Nested conditions required for selective field merging
	if cluster.Spec.Propagation.Spec != nil {
		clusterSpec := cluster.Spec.Propagation.Spec

		// Merge ingress configurations
		if clusterSpec.DelveIngress != nil {
			propagation.Spec.DelveIngress = clusterSpec.DelveIngress
		}
		if clusterSpec.QuicIngress != nil {
			propagation.Spec.QuicIngress = clusterSpec.QuicIngress
		}
		if clusterSpec.GrpcIngress != nil {
			propagation.Spec.GrpcIngress = clusterSpec.GrpcIngress
		}
		if clusterSpec.HTTPIngress != nil {
			propagation.Spec.HTTPIngress = clusterSpec.HTTPIngress
		}
		if clusterSpec.ProfilerIngress != nil {
			propagation.Spec.ProfilerIngress = clusterSpec.ProfilerIngress
		}
		if clusterSpec.ServiceAnnotations != nil {
			propagation.Spec.ServiceAnnotations = clusterSpec.ServiceAnnotations
		}

		// Merge deployment overrides selectively
		if clusterSpec.DeploymentOverrides != nil {
			if propagation.Spec.DeploymentOverrides == nil {
				propagation.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
			}
			mergeDeploymentOverrides(propagation.Spec.DeploymentOverrides, clusterSpec.DeploymentOverrides)
		}
	}

	// Apply cluster-level defaults
	if propagation.Spec.DeploymentOverrides == nil {
		propagation.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" && propagation.Spec.DeploymentOverrides.Image == "" {
		propagation.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}
	// Always apply cluster-level ImagePullSecrets (they override or are the default)
	if cluster.Spec.ImagePullSecrets != nil {
		propagation.Spec.DeploymentOverrides.ImagePullSecrets = cluster.Spec.ImagePullSecrets
	}

	return nil
}

func defaultPropagationSpec() *teranodev1alpha1.PropagationSpec {
	return &teranodev1alpha1.PropagationSpec{}
}
