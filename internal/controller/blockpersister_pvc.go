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
func (r *BlockPersisterReconciler) ReconcilePVC(log logr.Logger) (bool, error) {
	blockPersister := teranodev1alpha1.BlockPersister{}
	if err := r.Get(r.Context, r.NamespacedName, &blockPersister); err != nil {
		return false, err
	}
	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "block-persister-storage",
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}

	_, _ = controllerutil.CreateOrUpdate(r.Context, r.Client, &pvc, func() error {
		return r.updatePVC(&pvc, &blockPersister)
	})
	// for now ignore errors since there are immutable fields
	/*if err != nil && !k8serrors.IsForbidden(err) {
		return false, err
	}*/
	return true, nil
}

func (r *BlockPersisterReconciler) updatePVC(pvc *corev1.PersistentVolumeClaim, blockPersister *teranodev1alpha1.BlockPersister) error {
	err := controllerutil.SetControllerReference(blockPersister, pvc, r.Scheme)
	if err != nil {
		return err
	}
	pvc.Spec = *defaultBlockPersisterPVCSpec()

	// If storage class is configured, use it
	if blockPersister.Spec.StorageClass != "" {
		pvc.Spec.StorageClassName = &blockPersister.Spec.StorageClass
	}
	// If storage resources are configured, use them
	if blockPersister.Spec.StorageResources != nil {
		pvc.Spec.Resources = *blockPersister.Spec.StorageResources
	}
	return nil
}

// TODO: consider making it a generic functions, since it is a verbatim copy of the one from subtree validator
func defaultBlockPersisterPVCSpec() *corev1.PersistentVolumeClaimSpec {
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
