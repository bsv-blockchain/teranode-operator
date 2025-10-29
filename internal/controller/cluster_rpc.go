package controller

import (
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileRPC is the cluster rpc reconciler
func (r *ClusterReconciler) ReconcileRPC(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	rpc := teranodev1alpha1.RPC{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-rpc", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels("rpc"),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.RPC.Enabled || (cluster.Spec.Enabled != nil && !*cluster.Spec.Enabled) {
		namespacedName := types.NamespacedName{
			Name:      rpc.Name,
			Namespace: rpc.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &rpc)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &rpc)
		return true, err
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &rpc, func() error {
		return r.updateRPC(&rpc, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updateRPC(rpc *teranodev1alpha1.RPC, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, rpc, r.Scheme)
	if err != nil {
		return err
	}
	rpc.Spec = *defaultRPCSpec()

	// if user configures a config map name
	if cluster.Spec.RPC.Spec != nil {
		rpc.Spec = *cluster.Spec.RPC.Spec
	}
	if rpc.Spec.DeploymentOverrides == nil {
		rpc.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" && rpc.Spec.DeploymentOverrides.Image == "" {
		rpc.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}
	if cluster.Spec.ImagePullSecrets != nil {
		rpc.Spec.DeploymentOverrides.ImagePullSecrets = cluster.Spec.ImagePullSecrets
	}

	return nil
}

func defaultRPCSpec() *teranodev1alpha1.RPCSpec {
	return &teranodev1alpha1.RPCSpec{}
}
