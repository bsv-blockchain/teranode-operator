package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileGrpcService is the asset grpc service reconciler
func (r *AssetReconciler) ReconcileGrpcService(log logr.Logger) (bool, error) {
	asset := teranodev1alpha1.Asset{}
	if err := r.Get(r.Context, r.NamespacedName, &asset); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "asset-grpc",
			Namespace: r.NamespacedName.Namespace,
			Annotations: map[string]string{
				"prometheus.io/port":   "9091",
				"prometheus.io/scrape": "true",
			},
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateGrpcService(&svc, &asset)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *AssetReconciler) updateGrpcService(svc *corev1.Service, asset *teranodev1alpha1.Asset) error {
	err := controllerutil.SetControllerReference(asset, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultAssetGrpcServiceSpec()
	for k, v := range asset.Spec.ServiceAnnotations {
		svc.Annotations[k] = v
	}

	return nil
}

func defaultAssetGrpcServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "asset",
	}
	return &corev1.ServiceSpec{
		Selector: labels,
		Ports: []corev1.ServicePort{
			{
				Name:       "asset-grpc",
				Port:       int32(AssetGRPCPort),
				TargetPort: intstr.FromInt32(AssetGRPCPort),
				Protocol:   corev1.ProtocolTCP,
			},
		},
	}
}
