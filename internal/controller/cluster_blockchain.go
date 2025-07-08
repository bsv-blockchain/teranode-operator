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

// ReconcileBlockchain is the cluster blockchain reconciler
func (r *ClusterReconciler) ReconcileBlockchain(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	blockchain := teranodev1alpha1.Blockchain{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-blockchain", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels("blockchain"),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.Blockchain.Enabled || (cluster.Spec.Enabled != nil && !*cluster.Spec.Enabled) {
		namespacedName := types.NamespacedName{
			Name:      blockchain.Name,
			Namespace: blockchain.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &blockchain)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &blockchain)
		return true, err
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &blockchain, func() error {
		return r.updateBlockchain(&blockchain, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updateBlockchain(blockchain *teranodev1alpha1.Blockchain, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, blockchain, r.Scheme)
	if err != nil {
		return err
	}
	blockchain.Spec = *defaultBlockchainSpec()

	// if user configures a spec
	if cluster.Spec.Blockchain.Spec != nil {
		blockchain.Spec = *cluster.Spec.Blockchain.Spec
	}
	if blockchain.Spec.DeploymentOverrides == nil {
		blockchain.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" && blockchain.Spec.DeploymentOverrides.Image == "" {
		blockchain.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}

	return nil
}

func defaultBlockchainSpec() *teranodev1alpha1.BlockchainSpec {
	return &teranodev1alpha1.BlockchainSpec{}
}
