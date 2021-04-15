/*


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

package controllers

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	uselessv1alpha1 "github.com/NautiluX/useless-operator/api/v1alpha1"
)

// UselessObjectReconciler reconciles a UselessObject object
type UselessObjectReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=useless.useless.io,resources=uselessobjects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=useless.useless.io,resources=uselessobjects/status,verbs=get;update;patch

func (r *UselessObjectReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("uselessobject", req.NamespacedName)

	uselessCr := uselessv1alpha1.UselessObject{}
	err := r.Get(context.TODO(), req.NamespacedName, &uselessCr)
	if err != nil {
		return ctrl.Result{}, err
	}
	utc, _ := time.Now().UTC().MarshalText()

	uselessCr.Status.LastUpdatedAt = string(utc)

	err = r.Status().Update(context.TODO(), &uselessCr)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func ignoreStatusUpdatePredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Ignore updates to CR status in which case metadata.Generation does not change
			return e.MetaOld.GetGeneration() != e.MetaNew.GetGeneration()
		},
	}
}

func (r *UselessObjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&uselessv1alpha1.UselessObject{}).
		WithEventFilter(ignoreStatusUpdatePredicate()).
		Complete(r)
}
