/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
	"github.com/bsv-blockchain/teranode-operator/internal/utils"
)

// AssetReconciler reconciles a Asset object
type AssetReconciler struct {
	client.Client

	Scheme         *runtime.Scheme
	Log            logr.Logger
	NamespacedName types.NamespacedName
	Context        context.Context //nolint:containedctx // Required for reconciler pattern
}

//+kubebuilder:rbac:groups=teranode.bsvblockchain.org,resources=assets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=teranode.bsvblockchain.org,resources=assets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=teranode.bsvblockchain.org,resources=assets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Asset object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *AssetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	result := ctrl.Result{}
	r.Log = log.FromContext(ctx).WithValues("asset", req.NamespacedName)
	r.Context = ctx
	r.NamespacedName = req.NamespacedName
	asset := teranodev1alpha1.Asset{}
	if err := r.Get(ctx, req.NamespacedName, &asset); err != nil {
		r.Log.Error(err, "unable to fetch asset CR")
		return result, nil
	}

	_, err := utils.ReconcileBatch(r.Log,
		r.ReconcileDeployment,
		r.ReconcileService,
		r.ReconcileHTTPIngress,
		r.ReconcileHTTPSIngress,
	)

	// Update scale status (replicas and selector) from deployment
	deployment := &appsv1.Deployment{}
	if getErr := r.Get(ctx, types.NamespacedName{
		Name:      AssetDeploymentName,
		Namespace: asset.Namespace,
	}, deployment); getErr == nil {
		replicas, selector := utils.GetScaleStatusFromDeployment(deployment)
		asset.Status.Replicas = replicas
		asset.Status.Selector = selector
	}

	if err != nil {
		apimeta.SetStatusCondition(&asset.Status.Conditions,
			metav1.Condition{
				Type:    teranodev1alpha1.ConditionReconciled,
				Status:  metav1.ConditionFalse,
				Reason:  teranodev1alpha1.ReconciledReasonError,
				Message: err.Error(),
			},
		)
		_ = r.Client.Status().Update(ctx, &asset)
		// Since error is written on the status, let's log it and requeue
		// Returning error here is redundant
		r.Log.Error(err, "requeuing object for reconciliation")
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second}, nil
	} else {
		apimeta.SetStatusCondition(&asset.Status.Conditions,
			metav1.Condition{
				Type:    teranodev1alpha1.ConditionReconciled,
				Status:  metav1.ConditionTrue,
				Reason:  teranodev1alpha1.ReconciledReasonComplete,
				Message: teranodev1alpha1.ReconcileCompleteMessage,
			},
		)
	}

	err = r.Client.Status().Update(ctx, &asset)

	if asset.Spec.DeploymentOverrides != nil && asset.Spec.DeploymentOverrides.Replicas != nil {
		if asset.Status.Replicas != *asset.Spec.DeploymentOverrides.Replicas {
			r.Log.Info("requeuing to monitor replica status", "status", asset.Status.Replicas, "spec", asset.Spec.DeploymentOverrides.Replicas)
			return ctrl.Result{RequeueAfter: time.Second}, nil
		}
	}

	return ctrl.Result{Requeue: false, RequeueAfter: 0}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *AssetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&teranodev1alpha1.Asset{}).
		Owns(&appsv1.Deployment{}).
		Owns(&networkingv1.Ingress{}).
		Owns(&corev1.Service{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}
