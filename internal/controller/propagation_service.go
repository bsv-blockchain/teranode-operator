package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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
			Annotations: map[string]string{
				"prometheus.io/port":   "9091",
				"prometheus.io/scrape": "true",
			},
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
				Name:       "propagation-delve",
				Port:       int32(4041),
				TargetPort: intstr.FromInt(4040),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "propagation-grpc",
				Port:       int32(8084),
				TargetPort: intstr.FromInt(8084),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "propagation-quic",
				Port:       int32(8384),
				TargetPort: intstr.FromInt(8384),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "propagation-http",
				Port:       int32(8833),
				TargetPort: intstr.FromInt(8833),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "profiler", // TODO: shouldn't we call these x-profiler, where x is respective service name?
				Port:       int32(9091),
				TargetPort: intstr.FromInt(9091),
				Protocol:   corev1.ProtocolTCP,
			},
		},
	}
}
