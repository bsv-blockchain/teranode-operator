package controller

import (
	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func defaultNetworkPolicySpec() *networkingv1.NetworkPolicySpec {
	return &networkingv1.NetworkPolicySpec{
		PodSelector: metav1.LabelSelector{
			MatchLabels: getAppLabels(""),
		},
		PolicyTypes: []networkingv1.PolicyType{
			networkingv1.PolicyTypeEgress,
		},
		Egress: []networkingv1.NetworkPolicyEgressRule{
			{
				To: []networkingv1.NetworkPolicyPeer{
					{
						NamespaceSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "blockchain",
							},
						},
					},
				},
			},
		},
	}
}

// ReconcileNetworkPolicy defines the network policy reconciler
func (r *ClusterReconciler) ReconcileNetworkPolicy(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	labels := getAppLabels("")
	np := networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "teranode",
			Namespace: r.NamespacedName.Namespace,
			Labels:    labels,
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &np, func() error {
		return r.updateNetworkPolicy(&np, &cluster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updateNetworkPolicy(np *networkingv1.NetworkPolicy, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, np, r.Scheme)
	if err != nil {
		return err
	}
	np.Spec = *defaultNetworkPolicySpec()

	return nil
}
