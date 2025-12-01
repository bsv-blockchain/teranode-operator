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

// ReconcilePruner is the cluster pruner reconciler
func (r *ClusterReconciler) ReconcilePruner(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	pruner := teranodev1alpha1.Pruner{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-pruner", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels("pruner"),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.Pruner.Enabled || (cluster.Spec.Enabled != nil && !*cluster.Spec.Enabled) {
		namespacedName := types.NamespacedName{
			Name:      pruner.Name,
			Namespace: pruner.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &pruner)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &pruner)
		return true, err
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &pruner, func() error {
		return r.updatePruner(&pruner, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updatePruner(pruner *teranodev1alpha1.Pruner, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, pruner, r.Scheme)
	if err != nil {
		return err
	}
	pruner.Spec = *defaultPrunerSpec()

	// if user configures a spec
	if cluster.Spec.Pruner.Spec != nil {
		pruner.Spec = *cluster.Spec.Pruner.Spec
	}
	if pruner.Spec.DeploymentOverrides == nil {
		pruner.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" && pruner.Spec.DeploymentOverrides.Image == "" {
		pruner.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}
	if cluster.Spec.ImagePullSecrets != nil {
		pruner.Spec.DeploymentOverrides.ImagePullSecrets = cluster.Spec.ImagePullSecrets
	}

	return nil
}

func defaultPrunerSpec() *teranodev1alpha1.PrunerSpec {
	return &teranodev1alpha1.PrunerSpec{}
}
