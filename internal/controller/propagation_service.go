package controller

import (
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
	"github.com/bsv-blockchain/teranode-operator/internal/utils"
)

// ReconcileService is the propagation service reconciler
func (r *PropagationReconciler) ReconcileService(log logr.Logger) (bool, error) {
	propagation := teranodev1alpha1.Propagation{}
	if err := r.Get(r.Context, r.NamespacedName, &propagation); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "propagation",
			Namespace: r.NamespacedName.Namespace,
			Labels:    utils.GetPrometheusLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &propagation)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *PropagationReconciler) updateService(svc *corev1.Service, propagation *teranodev1alpha1.Propagation) error {
	err := controllerutil.SetControllerReference(propagation, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultPropagationServiceSpec()
	if svc.Annotations == nil {
		svc.Annotations = map[string]string{}
	}
	for k, v := range propagation.Spec.ServiceAnnotations {
		svc.Annotations[k] = v
	}
	return nil
}

func defaultPropagationServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "propagation",
	}
	ipFamily := corev1.IPFamilyPolicySingleStack
	return &corev1.ServiceSpec{
		Selector:       labels,
		ClusterIP:      "None",
		IPFamilyPolicy: &ipFamily,
		IPFamilies: []corev1.IPFamily{
			corev1.IPv4Protocol,
		},
		Ports: []corev1.ServicePort{
			{
				Name:       "propagation-grpc",
				Port:       int32(PropagationGRPCPort),
				TargetPort: intstr.FromInt32(PropagationGRPCPort),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "propagation-quic",
				Port:       int32(PropagationQuicPort),
				TargetPort: intstr.FromInt32(PropagationQuicPort),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "propagation-http",
				Port:       int32(PropagationHTTPPort),
				TargetPort: intstr.FromInt32(PropagationHTTPPort),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "health",
				Port:       int32(HealthPort),
				TargetPort: intstr.FromInt32(HealthPort),
				Protocol:   corev1.ProtocolTCP,
			},
		},
	}
}
