package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileHTTPIngress is the ingress for the asset grpc server
func (r *AssetReconciler) ReconcileHTTPIngress(log logr.Logger) (bool, error) {
	asset := teranodev1alpha1.Asset{}
	if err := r.Get(r.Context, r.NamespacedName, &asset); err != nil {
		return false, err
	}
	// Skip if HTTPIngress isn't set
	if asset.Spec.HTTPIngress == nil {
		return false, nil
	}
	labels := getAppLabels("asset")
	prefix := v1.PathTypePrefix
	ingress := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "asset-http",
			Namespace:   r.NamespacedName.Namespace,
			Annotations: asset.Spec.HTTPIngress.Annotations,
			Labels:      labels,
		},
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					Host: asset.Spec.GrpcIngress.Host,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &prefix,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: "asset",
											Port: v1.ServiceBackendPort{
												Name: "asset-http",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, ingress, func() error {
		return r.updateHTTPIngress(ingress, &asset)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *AssetReconciler) updateHTTPIngress(ingress *v1.Ingress, asset *teranodev1alpha1.Asset) error {
	err := controllerutil.SetControllerReference(asset, ingress, r.Scheme)
	if err != nil {
		return err
	}
	if asset.Spec.HTTPIngress == nil {
		return nil
	}
	if asset.Spec.HTTPIngress.Annotations == nil {
		ingress.Annotations = asset.Spec.HTTPIngress.Annotations
	}
	if asset.Spec.HTTPIngress.ClassName != nil {
		ingress.Spec.IngressClassName = asset.Spec.HTTPIngress.ClassName
	}

	return nil
}

// ReconcileHTTPSIngress is the ingress for the asset grpc server
func (r *AssetReconciler) ReconcileHTTPSIngress(log logr.Logger) (bool, error) {
	asset := teranodev1alpha1.Asset{}
	if err := r.Get(r.Context, r.NamespacedName, &asset); err != nil {
		return false, err
	}
	// Skip if domain isn't set
	if asset.Spec.HTTPSIngress == nil {
		return false, nil
	}

	labels := getAppLabels("asset")
	prefix := v1.PathTypePrefix
	ingress := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "asset-https",
			Namespace:   r.NamespacedName.Namespace,
			Annotations: asset.Spec.HTTPSIngress.Annotations,
			Labels:      labels,
		},
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					Host: asset.Spec.GrpcIngress.Host,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &prefix,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: "asset",
											Port: v1.ServiceBackendPort{
												Name: "asset-http",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			TLS: []v1.IngressTLS{
				{
					Hosts:      []string{asset.Spec.HTTPSIngress.Host},
					SecretName: "asset-tls",
				},
			},
		},
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, ingress, func() error {
		return r.updateHTTPSIngress(ingress, &asset)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *AssetReconciler) updateHTTPSIngress(ingress *v1.Ingress, asset *teranodev1alpha1.Asset) error {
	err := controllerutil.SetControllerReference(asset, ingress, r.Scheme)
	if err != nil {
		return err
	}
	if asset.Spec.HTTPSIngress == nil {
		return nil
	}
	if asset.Spec.HTTPSIngress.Annotations == nil {
		ingress.Annotations = asset.Spec.HTTPSIngress.Annotations
	}
	if asset.Spec.HTTPSIngress.ClassName != nil {
		ingress.Spec.IngressClassName = asset.Spec.HTTPSIngress.ClassName
	}

	return nil
}
