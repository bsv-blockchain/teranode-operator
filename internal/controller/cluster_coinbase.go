package controller

import (
	"fmt"

	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileCoinbase is the cluster coinbase reconciler
func (r *ClusterReconciler) ReconcileCoinbase(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	if !cluster.Spec.Coinbase.Enabled {
		return true, nil
	}
	coinbase := teranodev1alpha1.Coinbase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-coinbase", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &coinbase, func() error {
		return r.updateCoinbase(&coinbase, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updateCoinbase(coinbase *teranodev1alpha1.Coinbase, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, coinbase, r.Scheme)
	if err != nil {
		return err
	}
	coinbase.Spec = *defaultCoinbaseSpec()

	// if user configures a spec
	if cluster.Spec.Coinbase.Spec != nil {
		coinbase.Spec = *cluster.Spec.Coinbase.Spec
	}
	if coinbase.Spec.DeploymentOverrides == nil {
		coinbase.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" {
		coinbase.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}
	if cluster.Spec.ConfigMapName != "" {
		coinbase.Spec.DeploymentOverrides.ConfigMapName = cluster.Spec.ConfigMapName
	}
	return nil
}

func defaultCoinbaseSpec() *teranodev1alpha1.CoinbaseSpec {
	return &teranodev1alpha1.CoinbaseSpec{}
}
