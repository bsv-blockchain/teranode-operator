package controller

import (
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcilePropagation is the cluster propagation reconciler
func (r *ClusterReconciler) ReconcilePropagation(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	propagation := teranodev1alpha1.Propagation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-propagation", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels("propagation"),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.Propagation.Enabled {
		namespacedName := types.NamespacedName{
			Name:      propagation.Name,
			Namespace: propagation.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &propagation)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &propagation)
		return true, err
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &propagation, func() error {
		return r.updatePropagation(&propagation, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updatePropagation(propagation *teranodev1alpha1.Propagation, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, propagation, r.Scheme)
	if err != nil {
		return err
	}
	propagation.Spec = *defaultPropagationSpec()

	// if user configures a spec
	if cluster.Spec.Propagation.Spec != nil {
		propagation.Spec = *cluster.Spec.Propagation.Spec
	}
	if propagation.Spec.DeploymentOverrides == nil {
		propagation.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" && propagation.Spec.DeploymentOverrides.Image == "" {
		propagation.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}

	return nil
}

func defaultPropagationSpec() *teranodev1alpha1.PropagationSpec {
	return &teranodev1alpha1.PropagationSpec{}
}
