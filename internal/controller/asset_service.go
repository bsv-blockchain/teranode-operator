package controller

import (
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
	"github.com/bsv-blockchain/teranode-operator/internal/utils"
)

// ReconcileService is the asset service reconciler
func (r *AssetReconciler) ReconcileService(log logr.Logger) (bool, error) {
	asset := teranodev1alpha1.Asset{}
	if err := r.Get(r.Context, r.NamespacedName, &asset); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "asset",
			Namespace: r.NamespacedName.Namespace,
			Labels:    utils.GetPrometheusLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &asset)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *AssetReconciler) updateService(svc *corev1.Service, asset *teranodev1alpha1.Asset) error {
	err := controllerutil.SetControllerReference(asset, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultAssetServiceSpec()
	return nil
}

func defaultAssetServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "asset",
	}
	return &corev1.ServiceSpec{
		Selector: labels,
		Ports: []corev1.ServicePort{
			{
				Name:       "asset-http",
				Port:       int32(AssetHTTPPort),
				TargetPort: intstr.FromInt32(AssetHTTPPort),
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
