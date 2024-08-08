package controller

import (
	"fmt"
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileBlockAssembly is the node blockassembly reconciler
func (r *NodeReconciler) ReconcileBlockAssembly(log logr.Logger) (bool, error) {
	node := teranodev1alpha1.Node{}
	if err := r.Get(r.Context, r.NamespacedName, &node); err != nil {
		return false, err
	}
	if !node.Spec.BlockAssembly.Enabled {
		return true, nil
	}
	blockAssembly := teranodev1alpha1.BlockAssembly{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-blockassembly", node.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &blockAssembly, func() error {
		return r.updateBlockAssembly(&blockAssembly, &node)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *NodeReconciler) updateBlockAssembly(blockAssembly *teranodev1alpha1.BlockAssembly, node *teranodev1alpha1.Node) error {
	err := controllerutil.SetControllerReference(node, blockAssembly, r.Scheme)
	if err != nil {
		return err
	}
	blockAssembly.Spec = *defaultBlockAssemblySpec()

	// if user configures a config map name
	if node.Spec.BlockAssembly.Spec != nil {
		blockAssembly.Spec = *node.Spec.BlockAssembly.Spec
	}
	if node.Spec.Image != "" {
		blockAssembly.Spec.Image = node.Spec.Image
	}
	if node.Spec.ConfigMapName != "" {
		blockAssembly.Spec.ConfigMapName = node.Spec.ConfigMapName
	}
	return nil
}

func defaultBlockAssemblySpec() *teranodev1alpha1.BlockAssemblySpec {
	return &teranodev1alpha1.BlockAssemblySpec{}
}
