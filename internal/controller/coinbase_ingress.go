package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// TODO: Test that this actually works
// ReconcileGrpcIngress is the ingress for the coinbase grpc server
func (r *CoinbaseReconciler) ReconcileGrpcIngress(log logr.Logger) (bool, error) {

	coinbase := teranodev1alpha1.Coinbase{}
	if err := r.Get(r.Context, r.NamespacedName, &coinbase); err != nil {
		return false, err
	}
	// Skip if GrpcIngress isn't set
	if coinbase.Spec.GrpcIngress == nil {
		return false, nil
	}
	ingress := coinbase.Spec.GrpcIngress.DeepCopy()
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, ingress, func() error {
		return r.updateGrpcIngress(ingress, &coinbase)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *CoinbaseReconciler) updateGrpcIngress(ingress *v1.Ingress, coinbase *teranodev1alpha1.Coinbase) error {
	err := controllerutil.SetControllerReference(coinbase, ingress, r.Scheme)
	if err != nil {
		return err
	}
	ingress = coinbase.Spec.GrpcIngress.DeepCopy()
	return nil
}
