package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileService is the transaction blaster service reconciler
func (r *TransactionBlasterReconciler) ReconcileService(log logr.Logger) (bool, error) {
	transactionBlaster := teranodev1alpha1.TransactionBlaster{}
	if err := r.Get(r.Context, r.NamespacedName, &transactionBlaster); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "transaction-blaster",
			Namespace: r.NamespacedName.Namespace,
			Annotations: map[string]string{
				"prometheus.io/port":   "9092", // TODO: why isn't it 9091 as with other services?
				"prometheus.io/scrape": "true",
			},
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &transactionBlaster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *TransactionBlasterReconciler) updateService(svc *corev1.Service, transactionBlaster *teranodev1alpha1.TransactionBlaster) error {
	err := controllerutil.SetControllerReference(transactionBlaster, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultTransactionBlasterServiceSpec()
	return nil
}

func defaultTransactionBlasterServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "transactionBlaster",
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
				Name:       "transaction-blaster-tcp",
				Port:       int32(8086),
				TargetPort: intstr.FromInt(8086),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "profiler",
				Port:       int32(9091),
				TargetPort: intstr.FromInt(9091),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "transaction-blaster-profiler",
				Port:       int32(9092),
				TargetPort: intstr.FromInt(9092),
				Protocol:   corev1.ProtocolTCP,
			},
		},
	}
}
