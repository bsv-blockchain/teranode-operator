package controller

import (
	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
	"github.com/bsv-blockchain/teranode-operator/internal/utils"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileService is the bootstrap service reconciler
func (r *BootstrapReconciler) ReconcileService(log logr.Logger) (bool, error) {
	bs := teranodev1alpha1.Bootstrap{}
	if err := r.Get(r.Context, r.NamespacedName, &bs); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bootstrap",
			Namespace: r.NamespacedName.Namespace,
			Labels:    utils.GetPrometheusLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &bs)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BootstrapReconciler) updateService(svc *corev1.Service, bs *teranodev1alpha1.Bootstrap) error {
	err := controllerutil.SetControllerReference(bs, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultBootstrapServiceSpec()
	return nil
}

func defaultBootstrapServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "bootstrap",
	}
	return &corev1.ServiceSpec{
		Selector: labels,
		Ports: []corev1.ServicePort{
			{
				Name:       "bootstrap-http",
				Port:       int32(BootstrapHTTPPort),
				TargetPort: intstr.FromInt32(9906),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "bootstrap",
				Port:       int32(BootstrapGRPCPort),
				TargetPort: intstr.FromInt32(9905),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "delve",
				Port:       int32(DebuggerPort),
				TargetPort: intstr.FromInt32(4041),
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
