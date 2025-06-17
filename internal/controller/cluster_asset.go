package controller

import (
	"fmt"

	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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
	if !cluster.Spec.Asset.Enabled {
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

func (r *ClusterReconciler) updateAsset(asset *teranodev1alpha1.Asset, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, asset, r.Scheme)
	if err != nil {
		return err
	}
	asset.Spec = *defaultAssetSpec()

	// if user configures a spec
	if cluster.Spec.Asset.Spec != nil {
		asset.Spec = *cluster.Spec.Asset.Spec
	}
	if asset.Spec.DeploymentOverrides == nil {
		asset.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" && asset.Spec.DeploymentOverrides.Image == "" {
		asset.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}

	return nil
}

func defaultAssetSpec() *teranodev1alpha1.AssetSpec {
	return &teranodev1alpha1.AssetSpec{}
}
