package controller

import (
	"fmt"

	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileLegacy is the cluster coinbase reconciler
func (r *ClusterReconciler) ReconcileLegacy(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	if !cluster.Spec.Legacy.Enabled {
		return true, nil
	}
	legacy := teranodev1alpha1.Legacy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-legacy", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &legacy, func() error {
		return r.updateLegacy(&legacy, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updateLegacy(legacy *teranodev1alpha1.Legacy, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, legacy, r.Scheme)
	if err != nil {
		return err
	}
	legacy.Spec = *defaultLegacySpec()

	// if user configures a spec
	if cluster.Spec.Legacy.Spec != nil {
		legacy.Spec = *cluster.Spec.Legacy.Spec
	}
	if legacy.Spec.DeploymentOverrides == nil {
		legacy.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" {
		legacy.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}
	if cluster.Spec.ConfigMapName != "" {
		legacy.Spec.DeploymentOverrides.ConfigMapName = cluster.Spec.ConfigMapName
	}
	return nil
}

func defaultLegacySpec() *teranodev1alpha1.LegacySpec {
	return &teranodev1alpha1.LegacySpec{}
}
