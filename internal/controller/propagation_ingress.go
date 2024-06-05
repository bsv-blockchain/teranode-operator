package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// TODO: Test that this actually works

// ReconcileDelveIngress is the ingress for the propagation delve server
func (r *PropagationReconciler) ReconcileDelveIngress(log logr.Logger) (bool, error) {

	propagation := teranodev1alpha1.Propagation{}
	if err := r.Get(r.Context, r.NamespacedName, &propagation); err != nil {
		return false, err
	}
	// Skip if HttpIngress isn't set
	if propagation.Spec.DelveIngress == nil {
		return false, nil
	}
	ingress := propagation.Spec.DelveIngress.DeepCopy()
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, ingress, func() error {
		return r.updateDelveIngress(ingress, &propagation)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *PropagationReconciler) updateDelveIngress(ingress *v1.Ingress, propagation *teranodev1alpha1.Propagation) error {
	err := controllerutil.SetControllerReference(propagation, ingress, r.Scheme)
	if err != nil {
		return err
	}
	ingress = propagation.Spec.DelveIngress.DeepCopy()
	return nil
}

// ReconcileGrpcIngress is the ingress for the propagation grpc server
func (r *PropagationReconciler) ReconcileGrpcIngress(log logr.Logger) (bool, error) {

	propagation := teranodev1alpha1.Propagation{}
	if err := r.Get(r.Context, r.NamespacedName, &propagation); err != nil {
		return false, err
	}
	// Skip if GrpcIngress isn't set
	if propagation.Spec.GrpcIngress == nil {
		return false, nil
	}
	ingress := propagation.Spec.GrpcIngress.DeepCopy()
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, ingress, func() error {
		return r.updateGrpcIngress(ingress, &propagation)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *PropagationReconciler) updateGrpcIngress(ingress *v1.Ingress, propagation *teranodev1alpha1.Propagation) error {
	err := controllerutil.SetControllerReference(propagation, ingress, r.Scheme)
	if err != nil {
		return err
	}
	ingress = propagation.Spec.GrpcIngress.DeepCopy()
	return nil
}

// ReconcileDelveIngress is the ingress for the propagation delve server
func (r *PropagationReconciler) ReconcileQuicIngress(log logr.Logger) (bool, error) {

	propagation := teranodev1alpha1.Propagation{}
	if err := r.Get(r.Context, r.NamespacedName, &propagation); err != nil {
		return false, err
	}
	// Skip if HttpIngress isn't set
	if propagation.Spec.QuicIngress == nil {
		return false, nil
	}
	ingress := propagation.Spec.QuicIngress.DeepCopy()
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, ingress, func() error {
		return r.updateQuicIngress(ingress, &propagation)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *PropagationReconciler) updateQuicIngress(ingress *v1.Ingress, propagation *teranodev1alpha1.Propagation) error {
	err := controllerutil.SetControllerReference(propagation, ingress, r.Scheme)
	if err != nil {
		return err
	}
	ingress = propagation.Spec.QuicIngress.DeepCopy()
	return nil
}

// ReconcileHttpIngress is the ingress for the propagation grpc server
func (r *PropagationReconciler) ReconcileHttpIngress(log logr.Logger) (bool, error) {

	propagation := teranodev1alpha1.Propagation{}
	if err := r.Get(r.Context, r.NamespacedName, &propagation); err != nil {
		return false, err
	}
	// Skip if HttpIngress isn't set
	if propagation.Spec.HttpIngress == nil {
		return false, nil
	}
	ingress := propagation.Spec.HttpIngress.DeepCopy()
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, ingress, func() error {
		return r.updateHttpIngress(ingress, &propagation)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *PropagationReconciler) updateHttpIngress(ingress *v1.Ingress, propagation *teranodev1alpha1.Propagation) error {
	err := controllerutil.SetControllerReference(propagation, ingress, r.Scheme)
	if err != nil {
		return err
	}
	ingress = propagation.Spec.HttpIngress.DeepCopy()
	return nil
}

// ReconcileHttpsIngress is the ingress for the propagation grpc server
func (r *PropagationReconciler) ReconcileProfilerIngress(log logr.Logger) (bool, error) {

	propagation := teranodev1alpha1.Propagation{}
	if err := r.Get(r.Context, r.NamespacedName, &propagation); err != nil {
		return false, err
	}
	// Skip if domain isn't set
	if propagation.Spec.ProfilerIngress == nil {
		return false, nil
	}
	ingress := propagation.Spec.ProfilerIngress.DeepCopy()
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, ingress, func() error {
		return r.updateProfilerIngress(ingress, &propagation)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *PropagationReconciler) updateProfilerIngress(ingress *v1.Ingress, propagation *teranodev1alpha1.Propagation) error {
	err := controllerutil.SetControllerReference(propagation, ingress, r.Scheme)
	if err != nil {
		return err
	}
	ingress = propagation.Spec.ProfilerIngress.DeepCopy()
	return nil
}
