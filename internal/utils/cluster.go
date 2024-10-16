package utils

import (
	"context"

	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetClusterOwner(client client.Client, ctx context.Context, obj metav1.ObjectMeta) *teranodev1alpha1.Cluster {
	cluster := teranodev1alpha1.Cluster{}
	// Attempt to get the parent Cluster CR from owner refs
	ownerRefs := obj.GetOwnerReferences()
	for _, ownerRef := range ownerRefs {
		if ownerRef.Kind == "Cluster" {
			if err := client.Get(
				ctx,
				types.NamespacedName{
					Name:      ownerRef.Name,
					Namespace: obj.Namespace,
				}, &cluster); err != nil && !k8serrors.IsNotFound(err) {
				return nil
			}
		}
	}
	return &cluster
}
