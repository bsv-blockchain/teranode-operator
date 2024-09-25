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

// ReconcileBlockPersister is the cluster blockpersister reconciler
func (r *ClusterReconciler) ReconcileBlockPersister(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	blockPersister := teranodev1alpha1.BlockPersister{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-blockpersister", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.BlockPersister.Enabled {
		namespacedName := types.NamespacedName{
			Name:      blockPersister.Name,
			Namespace: blockPersister.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &blockPersister)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &blockPersister)
		return true, err
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &blockPersister, func() error {
		return r.updateBlockPersister(&blockPersister, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updateBlockPersister(blockPersister *teranodev1alpha1.BlockPersister, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, blockPersister, r.Scheme)
	if err != nil {
		return err
	}
	blockPersister.Spec = *defaultBlockPersisterSpec()

	// if user configures a spec
	if cluster.Spec.BlockPersister.Spec != nil {
		blockPersister.Spec = *cluster.Spec.BlockPersister.Spec
	}
	if blockPersister.Spec.DeploymentOverrides == nil {
		blockPersister.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" {
		blockPersister.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}
	if cluster.Spec.ConfigMapName != "" {
		blockPersister.Spec.DeploymentOverrides.ConfigMapName = cluster.Spec.ConfigMapName
	}
	return nil
}

func defaultBlockPersisterSpec() *teranodev1alpha1.BlockPersisterSpec {
	return &teranodev1alpha1.BlockPersisterSpec{}
}
