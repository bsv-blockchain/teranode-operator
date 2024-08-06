package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// TODO: Test that this actually works
// ReconcileGrpcIngress is the ingress for the asset grpc server
func (r *AssetReconciler) ReconcileGrpcIngress(log logr.Logger) (bool, error) {

	asset := teranodev1alpha1.Asset{}
	if err := r.Get(r.Context, r.NamespacedName, &asset); err != nil {
		return false, err
	}
	// Skip if GrpcIngress isn't set
	if asset.Spec.GrpcIngress == nil {
		return false, nil
	}
	labels := getAppLabels()
	prefix := v1.PathTypePrefix
	ingress := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "asset-grpc",
			Namespace:   r.NamespacedName.Namespace,
			Annotations: asset.Spec.GrpcIngress.Annotations,
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
											Name: "asset-grpc",
											Port: v1.ServiceBackendPort{
												Name: "asset-grpc",
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
		return r.updateGrpcIngress(ingress, &asset)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *AssetReconciler) updateGrpcIngress(ingress *v1.Ingress, asset *teranodev1alpha1.Asset) error {
	err := controllerutil.SetControllerReference(asset, ingress, r.Scheme)
	if err != nil {
		return err
	}
	if asset.Spec.GrpcIngress == nil {
		return nil
	}
	if asset.Spec.GrpcIngress.Annotations == nil {
		ingress.Annotations = asset.Spec.GrpcIngress.Annotations
	}
	if asset.Spec.GrpcIngress.ClassName != nil {
		ingress.Spec.IngressClassName = asset.Spec.GrpcIngress.ClassName
	}

	return nil
}

// ReconcileHttpIngress is the ingress for the asset grpc server
func (r *AssetReconciler) ReconcileHttpIngress(log logr.Logger) (bool, error) {

	asset := teranodev1alpha1.Asset{}
	if err := r.Get(r.Context, r.NamespacedName, &asset); err != nil {
		return false, err
	}
	// Skip if HttpIngress isn't set
	if asset.Spec.HttpIngress == nil {
		return false, nil
	}
	labels := getAppLabels()
	prefix := v1.PathTypePrefix
	ingress := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "asset-http",
			Namespace:   r.NamespacedName.Namespace,
			Annotations: asset.Spec.HttpIngress.Annotations,
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
		return r.updateHttpIngress(ingress, &asset)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *AssetReconciler) updateHttpIngress(ingress *v1.Ingress, asset *teranodev1alpha1.Asset) error {
	err := controllerutil.SetControllerReference(asset, ingress, r.Scheme)
	if err != nil {
		return err
	}
	if asset.Spec.HttpIngress == nil {
		return nil
	}
	if asset.Spec.HttpIngress.Annotations == nil {
		ingress.Annotations = asset.Spec.HttpIngress.Annotations
	}
	if asset.Spec.HttpIngress.ClassName != nil {
		ingress.Spec.IngressClassName = asset.Spec.HttpIngress.ClassName
	}

	return nil
}

// ReconcileHttpsIngress is the ingress for the asset grpc server
func (r *AssetReconciler) ReconcileHttpsIngress(log logr.Logger) (bool, error) {

	asset := teranodev1alpha1.Asset{}
	if err := r.Get(r.Context, r.NamespacedName, &asset); err != nil {
		return false, err
	}
	// Skip if domain isn't set
	if asset.Spec.HttpsIngress == nil {
		return false, nil
	}

	labels := getAppLabels()
	prefix := v1.PathTypePrefix
	ingress := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "asset-https",
			Namespace:   r.NamespacedName.Namespace,
			Annotations: asset.Spec.HttpsIngress.Annotations,
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
					Hosts:      []string{asset.Spec.HttpsIngress.Host},
					SecretName: "asset-tls",
				},
			},
		},
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, ingress, func() error {
		return r.updateHttpsIngress(ingress, &asset)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *AssetReconciler) updateHttpsIngress(ingress *v1.Ingress, asset *teranodev1alpha1.Asset) error {
	err := controllerutil.SetControllerReference(asset, ingress, r.Scheme)
	if err != nil {
		return err
	}
	if asset.Spec.HttpsIngress == nil {
		return nil
	}
	if asset.Spec.HttpsIngress.Annotations == nil {
		ingress.Annotations = asset.Spec.HttpsIngress.Annotations
	}
	if asset.Spec.HttpsIngress.ClassName != nil {
		ingress.Spec.IngressClassName = asset.Spec.HttpsIngress.ClassName
	}

	return nil
}
