package controller

import (
	"fmt"
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcilePeer is the node peer reconciler
func (r *NodeReconciler) ReconcilePeer(log logr.Logger) (bool, error) {
	node := teranodev1alpha1.Node{}
	if err := r.Get(r.Context, r.NamespacedName, &node); err != nil {
		return false, err
	}
	if !node.Spec.Peer.Enabled {
		return true, nil
	}
	peer := teranodev1alpha1.Peer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-peer", node.Name),
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &peer, func() error {
		return r.updatePeer(&peer, &node)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *NodeReconciler) updatePeer(peer *teranodev1alpha1.Peer, node *teranodev1alpha1.Node) error {
	err := controllerutil.SetControllerReference(node, peer, r.Scheme)
	if err != nil {
		return err
	}
	peer.Spec = *defaultPeerSpec()

	// if user configures a config map name
	if node.Spec.Peer.Spec != nil {
		peer.Spec = *node.Spec.Peer.Spec
	}
	if node.Spec.ConfigMapName != "" {
		peer.Spec.ConfigMapName = node.Spec.ConfigMapName
	}
	return nil
}

func defaultPeerSpec() *teranodev1alpha1.PeerSpec {
	return &teranodev1alpha1.PeerSpec{}
}
