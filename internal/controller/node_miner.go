package controller

import (
	"fmt"
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileMiner is the node coinbase reconciler
func (r *NodeReconciler) ReconcileMiner(log logr.Logger) (bool, error) {
	node := teranodev1alpha1.Node{}
	if err := r.Get(r.Context, r.NamespacedName, &node); err != nil {
		return false, err
	}
	if !node.Spec.Miner.Enabled {
		return true, nil
	}
	miner := teranodev1alpha1.Miner{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-miner", node.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &miner, func() error {
		return r.updateMiner(&miner, &node)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *NodeReconciler) updateMiner(miner *teranodev1alpha1.Miner, node *teranodev1alpha1.Node) error {
	err := controllerutil.SetControllerReference(node, miner, r.Scheme)
	if err != nil {
		return err
	}
	miner.Spec = *defaultMinerSpec()

	// if user configures a config map name
	if node.Spec.Miner.Spec != nil {
		miner.Spec = *node.Spec.Miner.Spec
	}
	if node.Spec.Image != "" {
		miner.Spec.Image = node.Spec.Image
	}
	if node.Spec.ConfigMapName != "" {
		miner.Spec.ConfigMapName = node.Spec.ConfigMapName
	}
	return nil
}

func defaultMinerSpec() *teranodev1alpha1.MinerSpec {
	return &teranodev1alpha1.MinerSpec{}
}
