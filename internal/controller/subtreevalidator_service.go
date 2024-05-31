package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileService is the subtree-validator service reconciler
func (r *SubtreeValidatorReconciler) ReconcileService(log logr.Logger) (bool, error) {
	subtreeValidator := teranodev1alpha1.SubtreeValidator{}
	if err := r.Get(r.Context, r.NamespacedName, &subtreeValidator); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "subtree-validator",
			Namespace: r.NamespacedName.Namespace,
			Annotations: map[string]string{
				"prometheus.io/port":   "9091",
				"prometheus.io/scrape": "true",
			},
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &subtreeValidator)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *SubtreeValidatorReconciler) updateService(svc *corev1.Service, subtreeValidator *teranodev1alpha1.SubtreeValidator) error {
	err := controllerutil.SetControllerReference(subtreeValidator, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultAssetServiceSpec()
	return nil
}

func defaultSubtreeValidatorServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "subtree-validator",
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
				Name:       "grpc",
				Port:       int32(8086),
				TargetPort: intstr.FromInt(8086),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "profiler",
				Port:       int32(9091),
				TargetPort: intstr.FromInt(9091),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "debugger",
				Port:       int32(4040),
				TargetPort: intstr.FromInt(4040),
				Protocol:   corev1.ProtocolTCP,
			},
		},
	}
}
