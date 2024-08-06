package controller

import (
	"fmt"
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileSubtreeValidator is the node subtreeValidator reconciler
func (r *NodeReconciler) ReconcileSubtreeValidator(log logr.Logger) (bool, error) {
	node := teranodev1alpha1.Node{}
	if err := r.Get(r.Context, r.NamespacedName, &node); err != nil {
		return false, err
	}
	if !node.Spec.SubtreeValidator.Enabled {
		return true, nil
	}
	subtreeValidator := teranodev1alpha1.SubtreeValidator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-subtreevalidator", node.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &subtreeValidator, func() error {
		return r.updateSubtreeValidator(&subtreeValidator, &node)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *NodeReconciler) updateSubtreeValidator(subtreeValidator *teranodev1alpha1.SubtreeValidator, node *teranodev1alpha1.Node) error {
	err := controllerutil.SetControllerReference(node, subtreeValidator, r.Scheme)
	if err != nil {
		return err
	}
	subtreeValidator.Spec = *defaultSubtreeValidatorSpec()

	// if user configures a config map name
	if node.Spec.SubtreeValidator.Spec != nil {
		subtreeValidator.Spec = *node.Spec.SubtreeValidator.Spec
	}
	if node.Spec.ConfigMapName != "" {
		subtreeValidator.Spec.ConfigMapName = node.Spec.ConfigMapName
	}
	return nil
}

func defaultSubtreeValidatorSpec() *teranodev1alpha1.SubtreeValidatorSpec {
	return &teranodev1alpha1.SubtreeValidatorSpec{}
}
