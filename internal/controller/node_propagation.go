package controller

import (
	"fmt"
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcilePropagation is the node propagation reconciler
func (r *NodeReconciler) ReconcilePropagation(log logr.Logger) (bool, error) {
	node := teranodev1alpha1.Node{}
	if err := r.Get(r.Context, r.NamespacedName, &node); err != nil {
		return false, err
	}
	if !node.Spec.Propagation.Enabled {
		return true, nil
	}
	propagation := teranodev1alpha1.Propagation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-propagation", node.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &propagation, func() error {
		return r.updatePropagation(&propagation, &node)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *NodeReconciler) updatePropagation(propagation *teranodev1alpha1.Propagation, node *teranodev1alpha1.Node) error {
	err := controllerutil.SetControllerReference(node, propagation, r.Scheme)
	if err != nil {
		return err
	}
	propagation.Spec = *defaultPropagationSpec()

	// if user configures a config map name
	if node.Spec.Propagation.Spec != nil {
		propagation.Spec = *node.Spec.Propagation.Spec
	}
	if node.Spec.Image != "" {
		propagation.Spec.Image = node.Spec.Image
	}
	if node.Spec.ConfigMapName != "" {
		propagation.Spec.ConfigMapName = node.Spec.ConfigMapName
	}
	return nil
}

func defaultPropagationSpec() *teranodev1alpha1.PropagationSpec {
	return &teranodev1alpha1.PropagationSpec{}
}
