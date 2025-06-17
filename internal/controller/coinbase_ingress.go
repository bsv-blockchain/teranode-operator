package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

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
	labels := getAppLabels("coinbase")
	prefix := v1.PathTypePrefix
	ingress := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "coinbase-grpc",
			Namespace:   r.NamespacedName.Namespace,
			Annotations: map[string]string{},
			Labels:      labels,
		},
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					Host: coinbase.Spec.GrpcIngress.Host,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &prefix,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: "coinbase",
											Port: v1.ServiceBackendPort{
												Name: "coinbase-tcp",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
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
	if coinbase.Spec.GrpcIngress == nil {
		return nil
	}
	if coinbase.Spec.GrpcIngress.Annotations == nil {
		ingress.Annotations = coinbase.Spec.GrpcIngress.Annotations
	}
	if coinbase.Spec.GrpcIngress.ClassName != nil {
		ingress.Spec.IngressClassName = coinbase.Spec.GrpcIngress.ClassName
	}
	return nil
}
