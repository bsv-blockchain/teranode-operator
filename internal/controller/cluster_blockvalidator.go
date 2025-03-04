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

// ReconcileBlockValidator is the cluster blockvalidator reconciler
func (r *ClusterReconciler) ReconcileBlockValidator(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	blockValidator := teranodev1alpha1.BlockValidator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-blockvalidator", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.BlockValidator.Enabled {
		namespacedName := types.NamespacedName{
			Name:      blockValidator.Name,
			Namespace: blockValidator.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &blockValidator)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &blockValidator)
		return true, err
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &blockValidator, func() error {
		return r.updateBlockValidator(&blockValidator, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updateBlockValidator(blockValidator *teranodev1alpha1.BlockValidator, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, blockValidator, r.Scheme)
	if err != nil {
		return err
	}
	blockValidator.Spec = *defaultBlockValidatorSpec()

	// if user configures a config map name
	if cluster.Spec.BlockValidator.Spec != nil {
		blockValidator.Spec = *cluster.Spec.BlockValidator.Spec
	}
	if blockValidator.Spec.DeploymentOverrides == nil {
		blockValidator.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" && blockValidator.Spec.DeploymentOverrides.Image == "" {
		blockValidator.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}

	return nil
}

func defaultBlockValidatorSpec() *teranodev1alpha1.BlockValidatorSpec {
	return &teranodev1alpha1.BlockValidatorSpec{}
}
