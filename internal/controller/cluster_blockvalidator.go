package controller

import (
	"fmt"

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
	if !cluster.Spec.BlockValidator.Enabled {
		return true, nil
	}
	blockValidator := teranodev1alpha1.BlockValidator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-blockvalidator", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
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
	if cluster.Spec.Image != "" {
		blockValidator.Spec.Image = cluster.Spec.Image
	}
	if cluster.Spec.ConfigMapName != "" {
		blockValidator.Spec.ConfigMapName = cluster.Spec.ConfigMapName
	}
	return nil
}

func defaultBlockValidatorSpec() *teranodev1alpha1.BlockValidatorSpec {
	return &teranodev1alpha1.BlockValidatorSpec{}
}
