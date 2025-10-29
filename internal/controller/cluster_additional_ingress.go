package controller

import (
	"fmt"

	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileAdditionalIngresses defines the network policy reconciler
func (r *ClusterReconciler) ReconcileAdditionalIngresses(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}

	labels := getAppLabels("cluster")
	for i, ingressSpec := range cluster.Spec.AdditionalIngresses {
		ingress := networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("teranode-%d", i),
				Namespace: r.NamespacedName.Namespace,
				Labels:    labels,
			},
		}
		_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &ingress, func() error {
			return r.updateIngress(&ingress, &ingressSpec, &cluster)
		})
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (r *ClusterReconciler) updateIngress(i *networkingv1.Ingress, spec *networkingv1.IngressSpec, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, i, r.Scheme)
	if err != nil {
		return err
	}
	i.Spec = *spec
	return nil
}
