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

// ReconcileService is the miner service reconciler
func (r *MinerReconciler) ReconcileService(log logr.Logger) (bool, error) {
	miner := teranodev1alpha1.Miner{}
	if err := r.Get(r.Context, r.NamespacedName, &miner); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "miner",
			Namespace: r.NamespacedName.Namespace,
			Labels:    utils.GetPrometheusLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &miner)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *MinerReconciler) updateService(svc *corev1.Service, miner *teranodev1alpha1.Miner) error {
	err := controllerutil.SetControllerReference(miner, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultMinerServiceSpec()
	return nil
}

func defaultMinerServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "miner",
	}
	return &corev1.ServiceSpec{
		Selector: labels,
		Ports: []corev1.ServicePort{
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
			{
				Name:       "health",
				Port:       int32(HealthPort),
				TargetPort: intstr.FromInt32(HealthPort),
				Protocol:   corev1.ProtocolTCP,
			},
		},
	}
}
