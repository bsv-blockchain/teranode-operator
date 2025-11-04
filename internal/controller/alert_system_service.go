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

// ReconcileService is the alert service reconciler
func (r *AlertSystemReconciler) ReconcileService(log logr.Logger) (bool, error) {
	alert := teranodev1alpha1.AlertSystem{}
	if err := r.Get(r.Context, r.NamespacedName, &alert); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "alert",
			Namespace: r.NamespacedName.Namespace,
			Labels:    utils.GetPrometheusLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &alert)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *AlertSystemReconciler) updateService(svc *corev1.Service, alert *teranodev1alpha1.AlertSystem) error {
	err := controllerutil.SetControllerReference(alert, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultAlertSystemServiceSpec()
	return nil
}

func defaultAlertSystemServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "alert",
	}
	return &corev1.ServiceSpec{
		Selector: labels,
		Ports: []corev1.ServicePort{
			{
				Name:       "alert-p2p",
				Port:       int32(AlertSystemPort),
				TargetPort: intstr.FromInt32(AlertSystemPort),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "health",
				Port:       int32(HealthPort),
				TargetPort: intstr.FromInt32(HealthPort),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "webserver",
				Port:       int32(AlertWebserverPort),
				TargetPort: intstr.FromInt32(AlertWebserverPort),
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
