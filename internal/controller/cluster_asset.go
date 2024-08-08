package controller

import (
	"fmt"
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileDeployment is the cluster service deployment reconciler
func (r *ClusterReconciler) ReconcileAsset(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	if !cluster.Spec.Asset.Enabled {
		return true, nil
	}
	asset := teranodev1alpha1.Asset{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-asset", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
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

	// if user configures a config map name
	if cluster.Spec.Asset.Spec != nil {
		asset.Spec = *cluster.Spec.Asset.Spec
	}
	if cluster.Spec.Image != "" {
		asset.Spec.Image = cluster.Spec.Image
	}
	if cluster.Spec.ConfigMapName != "" {
		asset.Spec.ConfigMapName = cluster.Spec.ConfigMapName
	}
	return nil
}

func defaultAssetSpec() *teranodev1alpha1.AssetSpec {
	return &teranodev1alpha1.AssetSpec{}
}
