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

// ReconcileService is the utxop service reconciler
func (r *UtxoPersisterReconciler) ReconcileService(log logr.Logger) (bool, error) {
	utxop := teranodev1alpha1.UtxoPersister{}
	if err := r.Get(r.Context, r.NamespacedName, &utxop); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "utxo-persister",
			Namespace: r.NamespacedName.Namespace,
			Labels:    utils.GetPrometheusLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &utxop)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UtxoPersisterReconciler) updateService(svc *corev1.Service, utxop *teranodev1alpha1.UtxoPersister) error {
	err := controllerutil.SetControllerReference(utxop, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultUtxoPersisterServiceSpec()
	return nil
}

func defaultUtxoPersisterServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "utxo-persister",
	}
	return &corev1.ServiceSpec{
		Selector: labels,
		Ports: []corev1.ServicePort{
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
