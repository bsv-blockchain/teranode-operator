package controller

import (
	"fmt"
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileBootstrap is the node bootstrap reconciler
func (r *NodeReconciler) ReconcileBootstrap(log logr.Logger) (bool, error) {
	node := teranodev1alpha1.Node{}
	if err := r.Get(r.Context, r.NamespacedName, &node); err != nil {
		return false, err
	}
	if !node.Spec.Bootstrap.Enabled {
		return true, nil
	}
	bootstrap := teranodev1alpha1.Bootstrap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-bootstrap", node.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &bootstrap, func() error {
		return r.updateBootstrap(&bootstrap, &node)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *NodeReconciler) updateBootstrap(bootstrap *teranodev1alpha1.Bootstrap, node *teranodev1alpha1.Node) error {
	err := controllerutil.SetControllerReference(node, bootstrap, r.Scheme)
	if err != nil {
		return err
	}
	bootstrap.Spec = *defaultBootstrapSpec()

	// if user configures a config map name
	if node.Spec.Bootstrap.Spec != nil {
		bootstrap.Spec = *node.Spec.Bootstrap.Spec
	}
	if node.Spec.ConfigMapName != "" {
		bootstrap.Spec.ConfigMapName = node.Spec.ConfigMapName
	}
	return nil
}

func defaultBootstrapSpec() *teranodev1alpha1.BootstrapSpec {
	return &teranodev1alpha1.BootstrapSpec{}
}
