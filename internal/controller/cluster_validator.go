package controller

import (
	"fmt"

	"github.com/go-logr/logr"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
)

// ReconcileValidator is the cluster validator reconciler
func (r *ClusterReconciler) ReconcileValidator(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	validator := teranodev1alpha1.Validator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-validator", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels("validator"),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.Validator.Enabled || (cluster.Spec.Enabled != nil && !*cluster.Spec.Enabled) {
		namespacedName := types.NamespacedName{
			Name:      validator.Name,
			Namespace: validator.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &validator)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &validator)
		return true, err
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

	// Selectively merge cluster spec - only override fields that are explicitly set
	if cluster.Spec.Validator.Spec != nil {
		validator.Spec = *cluster.Spec.Validator.Spec
	}

	if validator.Spec.DeploymentOverrides == nil {
		validator.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" {
		validator.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}
	if cluster.Spec.ImagePullSecrets != nil {
		validator.Spec.DeploymentOverrides.ImagePullSecrets = cluster.Spec.ImagePullSecrets
	}

	return nil
}

func defaultValidatorSpec() *teranodev1alpha1.ValidatorSpec {
	return &teranodev1alpha1.ValidatorSpec{}
}
