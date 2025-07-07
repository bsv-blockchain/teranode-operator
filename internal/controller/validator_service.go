package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/bitcoin-sv/teranode-operator/internal/utils"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileService is the validator service reconciler
func (r *ValidatorReconciler) ReconcileService(log logr.Logger) (bool, error) {
	validator := teranodev1alpha1.Validator{}
	if err := r.Get(r.Context, r.NamespacedName, &validator); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "validator",
			Namespace: r.NamespacedName.Namespace,
			Labels:    utils.GetPrometheusLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &validator)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ValidatorReconciler) updateService(svc *corev1.Service, validator *teranodev1alpha1.Validator) error {
	err := controllerutil.SetControllerReference(validator, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultValidatorServiceSpec()
	return nil
}

func defaultValidatorServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "validator",
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
				Name:       "validator-tcp",
				Port:       int32(8081),
				TargetPort: intstr.FromInt32(8081),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "validator", // TODO: seriously, these names
				Port:       int32(8088),
				TargetPort: intstr.FromInt32(8088),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "health",
				Port:       int32(HealthPort),
				TargetPort: intstr.FromInt32(HealthPort),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "profiler",
				Port:       int32(ProfilerPort),
				TargetPort: intstr.FromInt32(ProfilerPort),
				Protocol:   corev1.ProtocolTCP,
			},
		},
	}
}
