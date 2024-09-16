package controller

import (
	"fmt"

	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileValidator is the cluster validator reconciler
func (r *ClusterReconciler) ReconcileValidator(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	if !cluster.Spec.Validator.Enabled {
		return true, nil
	}
	validator := teranodev1alpha1.Validator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-validator", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &validator, func() error {
		return r.updateValidator(&validator, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updateValidator(validator *teranodev1alpha1.Validator, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, validator, r.Scheme)
	if err != nil {
		return err
	}
	validator.Spec = *defaultValidatorSpec()

	// if user configures a config map name
	if cluster.Spec.Validator.Spec != nil {
		validator.Spec = *cluster.Spec.Validator.Spec
	}
	if cluster.Spec.Image != "" {
		validator.Spec.Image = cluster.Spec.Image
	}
	if cluster.Spec.ConfigMapName != "" {
		validator.Spec.ConfigMapName = cluster.Spec.ConfigMapName
	}
	return nil
}

func defaultValidatorSpec() *teranodev1alpha1.ValidatorSpec {
	return &teranodev1alpha1.ValidatorSpec{}
}
