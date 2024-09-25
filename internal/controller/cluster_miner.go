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

// ReconcileMiner is the cluster coinbase reconciler
func (r *ClusterReconciler) ReconcileMiner(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	miner := teranodev1alpha1.Miner{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-miner", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.Miner.Enabled {
		namespacedName := types.NamespacedName{
			Name:      miner.Name,
			Namespace: miner.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &miner)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &miner)
		return true, err
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &miner, func() error {
		return r.updateMiner(&miner, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updateMiner(miner *teranodev1alpha1.Miner, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, miner, r.Scheme)
	if err != nil {
		return err
	}
	miner.Spec = *defaultMinerSpec()

	// if user configures a config map name
	if cluster.Spec.Miner.Spec != nil {
		miner.Spec = *cluster.Spec.Miner.Spec
	}
	if miner.Spec.DeploymentOverrides == nil {
		miner.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" {
		miner.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}
	if cluster.Spec.ConfigMapName != "" {
		miner.Spec.DeploymentOverrides.ConfigMapName = cluster.Spec.ConfigMapName
	}
	return nil
}

func defaultMinerSpec() *teranodev1alpha1.MinerSpec {
	return &teranodev1alpha1.MinerSpec{}
}
