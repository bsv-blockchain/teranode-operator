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

// ReconcileService is the peer service reconciler
func (r *PeerReconciler) ReconcileService(log logr.Logger) (bool, error) {
	peer := teranodev1alpha1.Peer{}
	if err := r.Get(r.Context, r.NamespacedName, &peer); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "peer",
			Namespace: r.NamespacedName.Namespace,
			Labels:    utils.GetPrometheusLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &peer)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *PeerReconciler) updateService(svc *corev1.Service, peer *teranodev1alpha1.Peer) error {
	err := controllerutil.SetControllerReference(peer, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultPeerServiceSpec()
	return nil
}

func defaultPeerServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "peer",
	}
	return &corev1.ServiceSpec{
		Selector: labels,
		Ports: []corev1.ServicePort{
			{
				Name:       "p2p-http",
				Port:       int32(9906),
				TargetPort: intstr.FromInt32(9906),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "p2p",
				Port:       int32(PeerPort),
				TargetPort: intstr.FromInt32(PeerPort),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "legacy",
				Port:       int32(PeerLegacyPort),
				TargetPort: intstr.FromInt32(PeerLegacyPort),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "delve",
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
			{
				Name:       "profiler",
				Port:       int32(ProfilerPort),
				TargetPort: intstr.FromInt32(ProfilerPort),
				Protocol:   corev1.ProtocolTCP,
			},
		},
	}
}
