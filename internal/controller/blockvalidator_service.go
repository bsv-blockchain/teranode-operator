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

// ReconcileService is the block-validator service reconciler
func (r *BlockValidatorReconciler) ReconcileService(log logr.Logger) (bool, error) {
	blockValidator := teranodev1alpha1.BlockValidator{}
	if err := r.Get(r.Context, r.NamespacedName, &blockValidator); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "block-validation",
			Namespace: r.NamespacedName.Namespace,
			Labels:    utils.GetPrometheusLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &blockValidator)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BlockValidatorReconciler) updateService(svc *corev1.Service, blockValidator *teranodev1alpha1.BlockValidator) error {
	err := controllerutil.SetControllerReference(blockValidator, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultBlockValidatorServiceSpec()
	return nil
}

func defaultBlockValidatorServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "block-validator",
	}
	return &corev1.ServiceSpec{
		Selector: labels,
		Ports: []corev1.ServicePort{
			{
				Name:       "tcp",
				Port:       int32(BlockValidationGRPCPort),
				TargetPort: intstr.FromInt32(BlockValidationHTTPPort),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "http",
				Port:       int32(BlockValidationHTTPPort),
				TargetPort: intstr.FromInt32(BlockValidationHTTPPort),
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
