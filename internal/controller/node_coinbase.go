package controller

import (
	"fmt"
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileCoinbase is the node coinbase reconciler
func (r *NodeReconciler) ReconcileCoinbase(log logr.Logger) (bool, error) {
	node := teranodev1alpha1.Node{}
	if err := r.Get(r.Context, r.NamespacedName, &node); err != nil {
		return false, err
	}
	if !node.Spec.Coinbase.Enabled {
		return true, nil
	}
	coinbase := teranodev1alpha1.Coinbase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-coinbase", node.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &coinbase, func() error {
		return r.updateCoinbase(&coinbase, &node)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *NodeReconciler) updateCoinbase(coinbase *teranodev1alpha1.Coinbase, node *teranodev1alpha1.Node) error {
	err := controllerutil.SetControllerReference(node, coinbase, r.Scheme)
	if err != nil {
		return err
	}
	coinbase.Spec = *defaultCoinbaseSpec()

	// if user configures a config map name
	if node.Spec.Coinbase.Spec != nil {
		coinbase.Spec = *node.Spec.Coinbase.Spec
	}
	if node.Spec.Image != "" {
		coinbase.Spec.Image = node.Spec.Image
	}
	if node.Spec.ConfigMapName != "" {
		coinbase.Spec.ConfigMapName = node.Spec.ConfigMapName
	}
	return nil
}

func defaultCoinbaseSpec() *teranodev1alpha1.CoinbaseSpec {
	return &teranodev1alpha1.CoinbaseSpec{}
}
