package controller

import (
	"fmt"
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileBlockAssembly is the cluster blockassembly reconciler
func (r *ClusterReconciler) ReconcileBlockAssembly(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	if !cluster.Spec.BlockAssembly.Enabled {
		return true, nil
	}
	blockAssembly := teranodev1alpha1.BlockAssembly{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-blockassembly", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
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

	// if user configures a config map name
	if cluster.Spec.BlockAssembly.Spec != nil {
		blockAssembly.Spec = *cluster.Spec.BlockAssembly.Spec
	}
	if cluster.Spec.Image != "" {
		blockAssembly.Spec.Image = cluster.Spec.Image
	}
	if cluster.Spec.ConfigMapName != "" {
		blockAssembly.Spec.ConfigMapName = cluster.Spec.ConfigMapName
	}
	return nil
}

func defaultBlockAssemblySpec() *teranodev1alpha1.BlockAssemblySpec {
	return &teranodev1alpha1.BlockAssemblySpec{}
}
