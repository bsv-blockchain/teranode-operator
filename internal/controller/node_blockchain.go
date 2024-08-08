package controller

import (
	"fmt"
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileBlockchain is the node blockchain reconciler
func (r *NodeReconciler) ReconcileBlockchain(log logr.Logger) (bool, error) {
	node := teranodev1alpha1.Node{}
	if err := r.Get(r.Context, r.NamespacedName, &node); err != nil {
		return false, err
	}
	if !node.Spec.Blockchain.Enabled {
		return true, nil
	}
	blockchain := teranodev1alpha1.Blockchain{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-blockchain", node.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &blockchain, func() error {
		return r.updateBlockchain(&blockchain, &node)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *NodeReconciler) updateBlockchain(blockchain *teranodev1alpha1.Blockchain, node *teranodev1alpha1.Node) error {
	err := controllerutil.SetControllerReference(node, blockchain, r.Scheme)
	if err != nil {
		return err
	}
	blockchain.Spec = *defaultBlockchainSpec()

	// if user configures a config map name
	if node.Spec.Blockchain.Spec != nil {
		blockchain.Spec = *node.Spec.Blockchain.Spec
	}
	if node.Spec.Image != "" {
		blockchain.Spec.Image = node.Spec.Image
	}
	if node.Spec.ConfigMapName != "" {
		blockchain.Spec.ConfigMapName = node.Spec.ConfigMapName
	}
	return nil
}

func defaultBlockchainSpec() *teranodev1alpha1.BlockchainSpec {
	return &teranodev1alpha1.BlockchainSpec{}
}
