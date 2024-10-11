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

// ReconcileService is the faucet service reconciler
func (r *FaucetReconciler) ReconcileService(log logr.Logger) (bool, error) {
	faucet := teranodev1alpha1.Faucet{}
	if err := r.Get(r.Context, r.NamespacedName, &faucet); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "faucet",
			Namespace: r.NamespacedName.Namespace,
			Labels:    utils.GetPrometheusLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &faucet)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *FaucetReconciler) updateService(svc *corev1.Service, faucet *teranodev1alpha1.Faucet) error {
	err := controllerutil.SetControllerReference(faucet, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultFaucetServiceSpec()
	return nil
}

func defaultFaucetServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "faucet",
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
				Name:       "faucet-tcp",
				Port:       int32(8097),
				TargetPort: intstr.FromInt32(8097),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "profiler",
				Port:       int32(9091),
				TargetPort: intstr.FromInt32(9091),
				Protocol:   corev1.ProtocolTCP,
			},
		},
	}
}
