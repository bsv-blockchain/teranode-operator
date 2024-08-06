package controller

import (
	"fmt"
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileValidator is the node validator reconciler
func (r *NodeReconciler) ReconcileValidator(log logr.Logger) (bool, error) {
	node := teranodev1alpha1.Node{}
	if err := r.Get(r.Context, r.NamespacedName, &node); err != nil {
		return false, err
	}
	if !node.Spec.Validator.Enabled {
		return true, nil
	}
	validator := teranodev1alpha1.Validator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-validator", node.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &validator, func() error {
		return r.updateValidator(&validator, &node)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *NodeReconciler) updateValidator(validator *teranodev1alpha1.Validator, node *teranodev1alpha1.Node) error {
	err := controllerutil.SetControllerReference(node, validator, r.Scheme)
	if err != nil {
		return err
	}
	validator.Spec = *defaultValidatorSpec()

	// if user configures a config map name
	if node.Spec.Validator.Spec != nil {
		validator.Spec = *node.Spec.Validator.Spec
	}
	if node.Spec.ConfigMapName != "" {
		validator.Spec.ConfigMapName = node.Spec.ConfigMapName
	}
	return nil
}

func defaultValidatorSpec() *teranodev1alpha1.ValidatorSpec {
	return &teranodev1alpha1.ValidatorSpec{}
}
