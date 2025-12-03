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

// ReconcileAsset is the asset reconciler
func (r *ClusterReconciler) ReconcileAsset(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	asset := teranodev1alpha1.Asset{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-asset", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels("asset"),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.Asset.Enabled || (cluster.Spec.Enabled != nil && !*cluster.Spec.Enabled) {
		namespacedName := types.NamespacedName{
			Name:      asset.Name,
			Namespace: asset.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &asset)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &asset)
		return true, err
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &asset, func() error {
		return r.updateAsset(&asset, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

//nolint:gocognit,gocyclo // Function complexity is inherent to handling multiple override cases
func (r *ClusterReconciler) updateAsset(asset *teranodev1alpha1.Asset, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, asset, r.Scheme)
	if err != nil {
		return err
	}

	// Only set defaults if this is a new resource (no spec configured yet)
	if asset.Spec.DeploymentOverrides == nil && cluster.Spec.Asset.Spec == nil {
		asset.Spec = *defaultAssetSpec()
	}

	// Selectively merge cluster spec - only override fields that are explicitly set
	//nolint:nestif // Nested conditions required for selective field merging
	if cluster.Spec.Asset.Spec != nil {
		clusterSpec := cluster.Spec.Asset.Spec

		// Merge ingress configurations
		if clusterSpec.GrpcIngress != nil {
			asset.Spec.GrpcIngress = clusterSpec.GrpcIngress
		}
		if clusterSpec.HTTPIngress != nil {
			asset.Spec.HTTPIngress = clusterSpec.HTTPIngress
		}
		if clusterSpec.HTTPSIngress != nil {
			asset.Spec.HTTPSIngress = clusterSpec.HTTPSIngress
		}

		// Merge deployment overrides selectively
		if clusterSpec.DeploymentOverrides != nil {
			if asset.Spec.DeploymentOverrides == nil {
				asset.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
			}
			mergeDeploymentOverrides(asset.Spec.DeploymentOverrides, clusterSpec.DeploymentOverrides)
		}
	}

	// Apply cluster-level defaults
	if asset.Spec.DeploymentOverrides == nil {
		asset.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" && asset.Spec.DeploymentOverrides.Image == "" {
		asset.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}
	// Always apply cluster-level ImagePullSecrets (they override or are the default)
	if cluster.Spec.ImagePullSecrets != nil {
		asset.Spec.DeploymentOverrides.ImagePullSecrets = cluster.Spec.ImagePullSecrets
	}

	return nil
}

func defaultAssetSpec() *teranodev1alpha1.AssetSpec {
	return &teranodev1alpha1.AssetSpec{}
}
