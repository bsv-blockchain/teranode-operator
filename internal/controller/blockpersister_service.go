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

// ReconcileService is the blockPersister service reconciler
func (r *BlockPersisterReconciler) ReconcileService(log logr.Logger) (bool, error) {
	blockPersister := teranodev1alpha1.BlockPersister{}
	if err := r.Get(r.Context, r.NamespacedName, &blockPersister); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "block-persister",
			Namespace: r.NamespacedName.Namespace,
			Labels:    utils.GetPrometheusLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &blockPersister)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BlockPersisterReconciler) updateService(svc *corev1.Service, blockPersister *teranodev1alpha1.BlockPersister) error {
	err := controllerutil.SetControllerReference(blockPersister, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultBlockPersisterServiceSpec()
	return nil
}

func defaultBlockPersisterServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "block-persister",
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
		},
	}
}
