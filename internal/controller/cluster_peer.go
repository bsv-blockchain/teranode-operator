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

// ReconcilePeer is the cluster peer reconciler
func (r *ClusterReconciler) ReconcilePeer(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	peer := teranodev1alpha1.Peer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-peer", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels("peer"),
		},
	}

	// Delete resource if we are disabling it
	if !cluster.Spec.Peer.Enabled || (cluster.Spec.Enabled != nil && !*cluster.Spec.Enabled) {
		namespacedName := types.NamespacedName{
			Name:      peer.Name,
			Namespace: peer.Namespace,
		}
		err := r.Get(r.Context, namespacedName, &peer)
		if k8serrors.IsNotFound(err) {
			return true, nil
		}
		// attempt to delete the resource
		err = r.Delete(r.Context, &peer)
		return true, err
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &peer, func() error {
		return r.updatePeer(&peer, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updatePeer(peer *teranodev1alpha1.Peer, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, peer, r.Scheme)
	if err != nil {
		return err
	}
	peer.Spec = *defaultPeerSpec()

	// if user configures a config map name
	if cluster.Spec.Peer.Spec != nil {
		peer.Spec = *cluster.Spec.Peer.Spec
	}
	if peer.Spec.DeploymentOverrides == nil {
		peer.Spec.DeploymentOverrides = &teranodev1alpha1.DeploymentOverrides{}
	}
	if cluster.Spec.Image != "" && peer.Spec.DeploymentOverrides.Image == "" {
		peer.Spec.DeploymentOverrides.Image = cluster.Spec.Image
	}

	return nil
}

func defaultPeerSpec() *teranodev1alpha1.PeerSpec {
	return &teranodev1alpha1.PeerSpec{}
}
