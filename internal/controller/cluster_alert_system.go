package controller

import (
	"fmt"

	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileAlertSystem is the alert system reconciler
func (r *ClusterReconciler) ReconcileAlertSystem(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	alertSystem := teranodev1alpha1.AlertSystem{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-alert-system", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.AlertSystem.Enabled {
		namespacedName := types.NamespacedName{
			Name:      alertSystem.Name,
			Namespace: alertSystem.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &alertSystem)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &alertSystem)
		return true, err
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &alertSystem, func() error {
		return r.updateAlertSystem(&alertSystem, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updateAlertSystem(alertSystem *teranodev1alpha1.AlertSystem, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, alertSystem, r.Scheme)
	if err != nil {
		return err
	}
	alertSystem.Spec = *defaultAlertSystemSpec()

	// if user configures a spec
	if cluster.Spec.AlertSystem.Spec != nil {
		alertSystem.Spec = *cluster.Spec.AlertSystem.Spec
	}
	if alertSystem.Spec.DeploymentOverrides == nil {
		alertSystem.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" {
		alertSystem.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}

	return nil
}

func defaultAlertSystemSpec() *teranodev1alpha1.AlertSystemSpec {
	return &teranodev1alpha1.AlertSystemSpec{}
}
