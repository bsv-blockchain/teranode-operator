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

// ReconcileUtxoPersister is the cluster coinbase reconciler
func (r *ClusterReconciler) ReconcileUtxoPersister(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	up := teranodev1alpha1.UtxoPersister{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-utxo-persister", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels("utxo-persister"),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.UtxoPersister.Enabled || (cluster.Spec.Enabled != nil && !*cluster.Spec.Enabled) {
		namespacedName := types.NamespacedName{
			Name:      up.Name,
			Namespace: up.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &up)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &up)
		return true, err
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &up, func() error {
		return r.updateUtxoPersister(&up, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updateUtxoPersister(up *teranodev1alpha1.UtxoPersister, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, up, r.Scheme)
	if err != nil {
		return err
	}
	up.Spec = *defaultUtxoPersisterSpec()

	// if user configures a config map name
	if cluster.Spec.UtxoPersister.Spec != nil {
		up.Spec = *cluster.Spec.UtxoPersister.Spec
	}
	if up.Spec.DeploymentOverrides == nil {
		up.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" && up.Spec.DeploymentOverrides.Image == "" {
		up.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}
	if cluster.Spec.ImagePullSecrets != nil {
		up.Spec.DeploymentOverrides.ImagePullSecrets = cluster.Spec.ImagePullSecrets
	}

	return nil
}

func defaultUtxoPersisterSpec() *teranodev1alpha1.UtxoPersisterSpec {
	return &teranodev1alpha1.UtxoPersisterSpec{}
}
