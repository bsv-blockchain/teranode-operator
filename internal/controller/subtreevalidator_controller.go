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

	"github.com/bitcoin-sv/teranode-operator/internal/utils"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
)

// SubtreeValidatorReconciler reconciles a SubtreeValidator object
type SubtreeValidatorReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	Log            logr.Logger
	NamespacedName types.NamespacedName
	Context        context.Context
}

//+kubebuilder:rbac:groups=teranode.bsvblockchain.org,resources=subtreevalidators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=teranode.bsvblockchain.org,resources=subtreevalidators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=teranode.bsvblockchain.org,resources=subtreevalidators/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SubtreeValidator object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *SubtreeValidatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	result := ctrl.Result{}
	r.Log = log.FromContext(ctx).WithValues("subtree-validator", req.NamespacedName)
	r.Context = ctx
	r.NamespacedName = req.NamespacedName
	subtreeValidator := teranodev1alpha1.SubtreeValidator{}
	if err := r.Get(ctx, req.NamespacedName, &subtreeValidator); err != nil {
		r.Log.Error(err, "unable to fetch asset CR")
		return result, nil
	}

	_, err := utils.ReconcileBatch(r.Log,
		r.ReconcileDeployment,
		r.ReconcileService,
		r.ReconcilePVC,
	)

	if err != nil {
		apimeta.SetStatusCondition(&subtreeValidator.Status.Conditions,
			metav1.Condition{
				Type:    teranodev1alpha1.ConditionReconciled,
				Status:  metav1.ConditionFalse,
				Reason:  teranodev1alpha1.ReconciledReasonError,
				Message: err.Error(),
			},
		)
		_ = r.Client.Status().Update(ctx, &subtreeValidator)
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second}, err
	} else {
		apimeta.SetStatusCondition(&subtreeValidator.Status.Conditions,
			metav1.Condition{
				Type:    teranodev1alpha1.ConditionReconciled,
				Status:  metav1.ConditionTrue,
				Reason:  teranodev1alpha1.ReconciledReasonComplete,
				Message: teranodev1alpha1.ReconcileCompleteMessage,
			},
		)
	}

	err = r.Client.Status().Update(ctx, &subtreeValidator)
	return ctrl.Result{Requeue: false, RequeueAfter: 0}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *SubtreeValidatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&teranodev1alpha1.SubtreeValidator{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Complete(r)
}
