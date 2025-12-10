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

// ReconcileBlockAssembly is the cluster blockassembly reconciler
func (r *ClusterReconciler) ReconcileBlockAssembly(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	blockAssembly := teranodev1alpha1.BlockAssembly{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-blockassembly", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels("block-assembly"),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.BlockAssembly.Enabled || (cluster.Spec.Enabled != nil && !*cluster.Spec.Enabled) {
		namespacedName := types.NamespacedName{
			Name:      blockAssembly.Name,
			Namespace: blockAssembly.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &blockAssembly)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &blockAssembly)
		return true, err
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &blockAssembly, func() error {
		return r.updateBlockAssembly(&blockAssembly, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updateBlockAssembly(blockAssembly *teranodev1alpha1.BlockAssembly, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, blockAssembly, r.Scheme)
	if err != nil {
		return err
	}
	blockAssembly.Spec = *defaultBlockAssemblySpec()

	// if user configures a spec
	if cluster.Spec.BlockAssembly.Spec != nil {
		blockAssembly.Spec = *cluster.Spec.BlockAssembly.Spec
	}
	if blockAssembly.Spec.DeploymentOverrides == nil {
		blockAssembly.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" {
		blockAssembly.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}
	if cluster.Spec.ImagePullSecrets != nil {
		blockAssembly.Spec.DeploymentOverrides.ImagePullSecrets = cluster.Spec.ImagePullSecrets
	}

	return nil
}

func defaultBlockAssemblySpec() *teranodev1alpha1.BlockAssemblySpec {
	return &teranodev1alpha1.BlockAssemblySpec{}
}
