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

	"github.com/bsv-blockchain/teranode-operator/internal/utils"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
)

// CoinbaseReconciler reconciles a Coinbase object
type CoinbaseReconciler struct {
	client.Client

	Scheme         *runtime.Scheme
	Log            logr.Logger
	NamespacedName types.NamespacedName
	//nolint:containedctx // Required for reconciler pattern
	Context context.Context
}

//+kubebuilder:rbac:groups=teranode.bsvblockchain.org,resources=coinbases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=teranode.bsvblockchain.org,resources=coinbases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=teranode.bsvblockchain.org,resources=coinbases/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Coinbase object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *CoinbaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	result := ctrl.Result{}
	r.Log = log.FromContext(ctx).WithValues("coinbase", req.NamespacedName)
	r.Context = ctx
	r.NamespacedName = req.NamespacedName
	coinbase := teranodev1alpha1.Coinbase{}
	if err := r.Get(ctx, req.NamespacedName, &coinbase); err != nil {
		r.Log.Error(err, "unable to fetch coinbase CR")
		return result, nil
	}

	_, err := utils.ReconcileBatch(r.Log,
		r.ReconcileDeployment,
		r.ReconcileService,
		r.ReconcileGrpcIngress,
	)

	if err != nil {
		apimeta.SetStatusCondition(&coinbase.Status.Conditions,
			metav1.Condition{
				Type:    teranodev1alpha1.ConditionReconciled,
				Status:  metav1.ConditionFalse,
				Reason:  teranodev1alpha1.ReconciledReasonError,
				Message: err.Error(),
			},
		)
		_ = r.Client.Status().Update(ctx, &coinbase)
		// Since error is written on the status, let's log it and requeue
		// Returning error here is redundant
		r.Log.Error(err, "requeuing object for reconciliation")
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second}, nil
	} else {
		apimeta.SetStatusCondition(&coinbase.Status.Conditions,
			metav1.Condition{
				Type:    teranodev1alpha1.ConditionReconciled,
				Status:  metav1.ConditionTrue,
				Reason:  teranodev1alpha1.ReconciledReasonComplete,
				Message: teranodev1alpha1.ReconcileCompleteMessage,
			},
		)
	}

	err = r.Client.Status().Update(ctx, &coinbase)
	return ctrl.Result{Requeue: false, RequeueAfter: 0}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *CoinbaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&teranodev1alpha1.Coinbase{}).
		Owns(&appsv1.Deployment{}).
		Owns(&networkingv1.Ingress{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
