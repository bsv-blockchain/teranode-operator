package controller

import (
	"fmt"

	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileSubtreeValidator is the cluster subtreeValidator reconciler
func (r *ClusterReconciler) ReconcileSubtreeValidator(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	if !cluster.Spec.SubtreeValidator.Enabled {
		return true, nil
	}
	subtreeValidator := teranodev1alpha1.SubtreeValidator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-subtreevalidator", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &subtreeValidator, func() error {
		return r.updateSubtreeValidator(&subtreeValidator, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updateSubtreeValidator(subtreeValidator *teranodev1alpha1.SubtreeValidator, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, subtreeValidator, r.Scheme)
	if err != nil {
		return err
	}
	subtreeValidator.Spec = *defaultSubtreeValidatorSpec()

	// if user configures a config map name
	if cluster.Spec.SubtreeValidator.Spec != nil {
		subtreeValidator.Spec = *cluster.Spec.SubtreeValidator.Spec
	}
	if subtreeValidator.Spec.DeploymentOverrides == nil {
		subtreeValidator.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" {
		subtreeValidator.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}
	if cluster.Spec.ConfigMapName != "" {
		subtreeValidator.Spec.DeploymentOverrides.ConfigMapName = cluster.Spec.ConfigMapName
	}
	return nil
}

func defaultSubtreeValidatorSpec() *teranodev1alpha1.SubtreeValidatorSpec {
	return &teranodev1alpha1.SubtreeValidatorSpec{}
}
