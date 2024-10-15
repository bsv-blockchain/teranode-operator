package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcilePVC is the postgres PVC
func (r *SubtreeValidatorReconciler) ReconcilePVC(log logr.Logger) (bool, error) {
	subtreeValidator := teranodev1alpha1.SubtreeValidator{}
	if err := r.Get(r.Context, r.NamespacedName, &subtreeValidator); err != nil {
		return false, err
	}
	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      SharedPVCName,
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}

	_, _ = controllerutil.CreateOrUpdate(r.Context, r.Client, &pvc, func() error {
		return r.updatePVC(&pvc, &subtreeValidator)
	})
	/*
		// Ignore this check for right now until we properly implement resizing
		if err != nil && !k8serrors.IsForbidden(err) {
			return false, err
		}*/
	return true, nil
}

func (r *SubtreeValidatorReconciler) updatePVC(pvc *corev1.PersistentVolumeClaim, subtreeValidator *teranodev1alpha1.SubtreeValidator) error {
	err := controllerutil.SetControllerReference(subtreeValidator, pvc, r.Scheme)
	if err != nil {
		return err
	}
	pvc.Spec = *defaultPVCSpec()

	// If storage class is configured, use it
	if subtreeValidator.Spec.StorageClass != "" {
		pvc.Spec.StorageClassName = &subtreeValidator.Spec.StorageClass
	}
	// If storage resources are configured, use them
	if subtreeValidator.Spec.StorageResources != nil {
		pvc.Spec.Resources = *subtreeValidator.Spec.StorageResources
	}
	if subtreeValidator.Spec.StorageVolume != "" {
		pvc.Spec.VolumeName = subtreeValidator.Spec.StorageVolume
	}
	return nil
}

func defaultPVCSpec() *corev1.PersistentVolumeClaimSpec {
	emptyStorageClass := ""
	return &corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteMany,
		},
		StorageClassName: &emptyStorageClass,
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				"storage": resource.MustParse("2400Gi"),
			},
		},
	}
}
