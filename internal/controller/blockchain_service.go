package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileService is the blockchain service reconciler
func (r *BlockchainReconciler) ReconcileService(log logr.Logger) (bool, error) {
	blockchain := teranodev1alpha1.Blockchain{}
	if err := r.Get(r.Context, r.NamespacedName, &blockchain); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "blockchain",
			Namespace: r.NamespacedName.Namespace,
			Annotations: map[string]string{
				"prometheus.io/port":   "9091",
				"prometheus.io/scrape": "true",
			},
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &blockchain)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BlockchainReconciler) updateService(svc *corev1.Service, blockchain *teranodev1alpha1.Blockchain) error {
	err := controllerutil.SetControllerReference(blockchain, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultBlockchainServiceSpec()
	return nil
}

func defaultBlockchainServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "blockchain",
	}
	return &corev1.ServiceSpec{
		Selector: labels,
		Ports: []corev1.ServicePort{
			{
				Name:       "tcp",
				Port:       int32(8081),
				TargetPort: intstr.FromInt32(8081),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "blockchain",
				Port:       int32(8087),
				TargetPort: intstr.FromInt32(8087),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "profiler",
				Port:       int32(9091),
				TargetPort: intstr.FromInt32(9091),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "debugger",
				Port:       int32(4040),
				TargetPort: intstr.FromInt32(4040),
				Protocol:   corev1.ProtocolTCP,
			},
		},
	}
}
