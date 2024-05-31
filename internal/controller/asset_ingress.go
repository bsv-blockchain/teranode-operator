package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/networking/v1"
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
	ingress := asset.Spec.GrpcIngress.DeepCopy()
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
	ingress = asset.Spec.GrpcIngress.DeepCopy()
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
	ingress := asset.Spec.HttpIngress.DeepCopy()
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
	ingress = asset.Spec.HttpIngress.DeepCopy()
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
	ingress := asset.Spec.HttpsIngress.DeepCopy()
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
	ingress = asset.Spec.GrpcIngress.DeepCopy()
	return nil
}
