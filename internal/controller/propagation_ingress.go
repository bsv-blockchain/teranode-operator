package controller

import (
	"github.com/go-logr/logr"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
)

// ReconcileGrpcIngress is the ingress for the propagation grpc server
func (r *PropagationReconciler) ReconcileGrpcIngress(log logr.Logger) (bool, error) {
	propagation := teranodev1alpha1.Propagation{}
	if err := r.Get(r.Context, r.NamespacedName, &propagation); err != nil {
		return false, err
	}
	// Skip if GrpcIngress isn't set
	if propagation.Spec.GrpcIngress == nil {
		return false, nil
	}

	labels := getAppLabels("propagation")
	prefix := v1.PathTypePrefix
	ingress := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "propagation-grpc",
			Namespace:   r.NamespacedName.Namespace,
			Annotations: propagation.Spec.GrpcIngress.Annotations,
			Labels:      labels,
		},
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					Host: propagation.Spec.GrpcIngress.Host,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &prefix,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: "propagation",
											Port: v1.ServiceBackendPort{
												Number: PropagationGRPCPort,
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
		return r.updateGrpcIngress(ingress, &propagation)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *PropagationReconciler) updateGrpcIngress(ingress *v1.Ingress, propagation *teranodev1alpha1.Propagation) error {
	err := controllerutil.SetControllerReference(propagation, ingress, r.Scheme)
	if err != nil {
		return err
	}
	if propagation.Spec.GrpcIngress == nil {
		return nil
	}

	if propagation.Spec.GrpcIngress.Annotations == nil {
		ingress.Annotations = propagation.Spec.GrpcIngress.Annotations
	}
	if propagation.Spec.GrpcIngress.ClassName != nil {
		ingress.Spec.IngressClassName = propagation.Spec.GrpcIngress.ClassName
	}
	return nil
}
