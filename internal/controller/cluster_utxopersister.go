package controller

import (
	"fmt"

	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileUtxoPersister is the cluster coinbase reconciler
func (r *ClusterReconciler) ReconcileUtxoPersister(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	if !cluster.Spec.UtxoPersister.Enabled {
		return true, nil
	}
	up := teranodev1alpha1.UtxoPersister{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-utxo-persister", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
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
	if cluster.Spec.Image != "" {
		up.Spec.Image = cluster.Spec.Image
	}
	if cluster.Spec.ConfigMapName != "" {
		up.Spec.ConfigMapName = cluster.Spec.ConfigMapName
	}
	return nil
}

func defaultUtxoPersisterSpec() *teranodev1alpha1.UtxoPersisterSpec {
	return &teranodev1alpha1.UtxoPersisterSpec{}
}
