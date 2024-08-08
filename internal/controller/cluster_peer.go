package controller

import (
	"fmt"
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
	if !cluster.Spec.Peer.Enabled {
		return true, nil
	}
	peer := teranodev1alpha1.Peer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-peer", cluster.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
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
	if cluster.Spec.Image != "" {
		peer.Spec.Image = cluster.Spec.Image
	}
	if cluster.Spec.ConfigMapName != "" {
		peer.Spec.ConfigMapName = cluster.Spec.ConfigMapName
	}
	return nil
}

func defaultPeerSpec() *teranodev1alpha1.PeerSpec {
	return &teranodev1alpha1.PeerSpec{}
}
