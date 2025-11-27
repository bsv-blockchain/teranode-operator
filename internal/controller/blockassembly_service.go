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

// ReconcileService is the blockassembly service reconciler
func (r *BlockAssemblyReconciler) ReconcileService(log logr.Logger) (bool, error) {
	blockassembly := teranodev1alpha1.BlockAssembly{}
	if err := r.Get(r.Context, r.NamespacedName, &blockassembly); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "block-assembly",
			Namespace: r.NamespacedName.Namespace,
			Labels:    utils.GetPrometheusLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &blockassembly)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BlockAssemblyReconciler) updateService(svc *corev1.Service, blockassembly *teranodev1alpha1.BlockAssembly) error {
	err := controllerutil.SetControllerReference(blockassembly, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultBlockAssemblyServiceSpec()
	return nil
}

func defaultBlockAssemblyServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "block-assembly",
	}
	return &corev1.ServiceSpec{
		Selector: labels,
		Ports: []corev1.ServicePort{
			{
				Name:       "block-assembly",
				Port:       int32(BlockAssemblyPort),
				TargetPort: intstr.FromInt32(BlockAssemblyPort),
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
