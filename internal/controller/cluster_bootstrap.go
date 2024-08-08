package controller

import (
	"fmt"
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileBootstrap is the cluster bootstrap reconciler
func (r *ClusterReconciler) ReconcileBootstrap(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	if !cluster.Spec.Bootstrap.Enabled {
		return true, nil
	}
	bootstrap := teranodev1alpha1.Bootstrap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-bootstrap", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &bootstrap, func() error {
		return r.updateBootstrap(&bootstrap, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updateBootstrap(bootstrap *teranodev1alpha1.Bootstrap, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, bootstrap, r.Scheme)
	if err != nil {
		return err
	}
	bootstrap.Spec = *defaultBootstrapSpec()

	// if user configures a config map name
	if cluster.Spec.Bootstrap.Spec != nil {
		bootstrap.Spec = *cluster.Spec.Bootstrap.Spec
	}
	if cluster.Spec.Image != "" {
		bootstrap.Spec.Image = cluster.Spec.Image
	}
	if cluster.Spec.ConfigMapName != "" {
		bootstrap.Spec.ConfigMapName = cluster.Spec.ConfigMapName
	}
	return nil
}

func defaultBootstrapSpec() *teranodev1alpha1.BootstrapSpec {
	return &teranodev1alpha1.BootstrapSpec{}
}
