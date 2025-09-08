package controller

import (
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

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
	legacy := teranodev1alpha1.Legacy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-legacy", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels("legacy"),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.Legacy.Enabled || (cluster.Spec.Enabled != nil && !*cluster.Spec.Enabled) {
		namespacedName := types.NamespacedName{
			Name:      legacy.Name,
			Namespace: legacy.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &legacy)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &legacy)
		return true, err
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
	if cluster.Spec.Image != "" && legacy.Spec.DeploymentOverrides.Image == "" {
		legacy.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}
	if cluster.Spec.ImagePullSecrets != nil {
		legacy.Spec.DeploymentOverrides.ImagePullSecrets = cluster.Spec.ImagePullSecrets
	}

	return nil
}

func defaultLegacySpec() *teranodev1alpha1.LegacySpec {
	return &teranodev1alpha1.LegacySpec{}
}
