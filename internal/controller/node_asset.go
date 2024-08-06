package controller

import (
	"fmt"
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileDeployment is the node service deployment reconciler
func (r *NodeReconciler) ReconcileAsset(log logr.Logger) (bool, error) {
	node := teranodev1alpha1.Node{}
	if err := r.Get(r.Context, r.NamespacedName, &node); err != nil {
		return false, err
	}
	if !node.Spec.Asset.Enabled {
		return true, nil
	}
	asset := teranodev1alpha1.Asset{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-asset", node.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &asset, func() error {
		return r.updateAsset(&asset, &node)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *NodeReconciler) updateAsset(asset *teranodev1alpha1.Asset, node *teranodev1alpha1.Node) error {
	err := controllerutil.SetControllerReference(node, asset, r.Scheme)
	if err != nil {
		return err
	}
	asset.Spec = *defaultAssetSpec()

	// if user configures a config map name
	if node.Spec.Asset.Spec != nil {
		asset.Spec = *node.Spec.Asset.Spec
	}
	if node.Spec.ConfigMapName != "" {
		asset.Spec.ConfigMapName = node.Spec.ConfigMapName
	}
	return nil
}

func defaultAssetSpec() *teranodev1alpha1.AssetSpec {
	return &teranodev1alpha1.AssetSpec{}
}
