package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileService is the legacy service reconciler
func (r *LegacyReconciler) ReconcileService(log logr.Logger) (bool, error) {
	legacy := teranodev1alpha1.Legacy{}
	if err := r.Get(r.Context, r.NamespacedName, &legacy); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "legacy",
			Namespace: r.NamespacedName.Namespace,
			Annotations: map[string]string{
				"prometheus.io/port":   "9091",
				"prometheus.io/scrape": "true",
			},
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &legacy)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *LegacyReconciler) updateService(svc *corev1.Service, legacy *teranodev1alpha1.Legacy) error {
	err := controllerutil.SetControllerReference(legacy, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultLegacyServiceSpec()
	return nil
}

func defaultLegacyServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "legacy",
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
				Name:       "http",
				Port:       int32(LegacyHttpPort),
				TargetPort: intstr.FromInt32(LegacyHttpPort),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "profiler",
				Port:       int32(ProfilerPort),
				TargetPort: intstr.FromInt32(ProfilerPort),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "debugger",
				Port:       int32(DebuggerPort),
				TargetPort: intstr.FromInt32(DebuggerPort),
				Protocol:   corev1.ProtocolTCP,
			},
		},
	}
}
