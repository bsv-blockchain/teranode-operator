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

// ReconcileService is the rpc service reconciler
func (r *RPCReconciler) ReconcileService(log logr.Logger) (bool, error) {
	rpc := teranodev1alpha1.RPC{}
	if err := r.Get(r.Context, r.NamespacedName, &rpc); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rpc",
			Namespace: r.NamespacedName.Namespace,
			Labels:    utils.GetPrometheusLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &rpc)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RPCReconciler) updateService(svc *corev1.Service, rpc *teranodev1alpha1.RPC) error {
	err := controllerutil.SetControllerReference(rpc, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultRPCServiceSpec()
	return nil
}

func defaultRPCServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "rpc",
	}
	return &corev1.ServiceSpec{
		Selector: labels,
		Ports: []corev1.ServicePort{
			{
				Name:       "rpc",
				Port:       int32(RPCPort),
				TargetPort: intstr.FromInt32(RPCPort),
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
