package controller

import (
	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcilePVC is the postgres PVC
func (r *ClusterReconciler) ReconcilePVC(log logr.Logger) (bool, error) {
	cluster := teranodev1alpha1.Cluster{}
	if err := r.Get(r.Context, r.NamespacedName, &cluster); err != nil {
		return false, err
	}
	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      SharedPVCName,
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels("cluster"),
		},
	}
	// Check if PVC is already created so that we can copy the existing spec values
	// This is how we properly support resizing
	existingPvcNamespacedName := types.NamespacedName{
		Namespace: pvc.Namespace,
		Name:      pvc.Name,
	}
	existingPVC := &corev1.PersistentVolumeClaim{}
	err := r.Get(r.Context, existingPvcNamespacedName, existingPVC)
	if err != nil && !k8serrors.IsNotFound(err) {
		return false, err
	}

	// If in cluster PVC is not found; nil it out to not confuse the create or update section
	if k8serrors.IsNotFound(err) {
		existingPVC = nil
	}

	_, err = controllerutil.CreateOrUpdate(r.Context, r.Client, &pvc, func() error {
		return r.updatePVC(&pvc, existingPVC, &cluster)
	})

	// Ignore forbidden errors
	if err != nil && !k8serrors.IsForbidden(err) {
		return false, err
	}
	return true, nil
}

func (r *ClusterReconciler) updatePVC(pvc, inClusterPVC *corev1.PersistentVolumeClaim, cluster *teranodev1alpha1.Cluster) error {
	err := controllerutil.SetControllerReference(cluster, pvc, r.Scheme)
	if err != nil {
		return err
	}
	if inClusterPVC == nil {
		pvc.Spec = *defaultPVCSpec()
	} else {
		pvc.Spec = *inClusterPVC.Spec.DeepCopy()
	}

	// If storage class is configured, use it
	if cluster.Spec.SharedStorage.StorageClass != "" {
		pvc.Spec.StorageClassName = &cluster.Spec.SharedStorage.StorageClass
	}
	// If storage resources are configured, use them
	if cluster.Spec.SharedStorage.StorageResources != nil {
		pvc.Spec.Resources = *cluster.Spec.SharedStorage.StorageResources
	}
	if cluster.Spec.SharedStorage.StorageVolume != "" {
		pvc.Spec.VolumeName = cluster.Spec.SharedStorage.StorageVolume
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
		Resources: corev1.VolumeResourceRequirements{
			Requests: corev1.ResourceList{
				"storage": resource.MustParse("2400Gi"),
			},
		},
	}
}
