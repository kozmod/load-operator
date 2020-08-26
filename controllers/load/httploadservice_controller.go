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

	"github.com/kozmod/load-operator/domain/metrics/usecase"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	loadv1alpha1 "github.com/kozmod/load-operator/apis/load/v1alpha1"
)

// HttpLoadServiceReconciler reconciles a HttpLoadService object
type HttpLoadServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=load.load-operator.com,resources=httploadservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=load.load-operator.com,resources=httploadservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

func (r *HttpLoadServiceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	loadservice := &loadv1alpha1.HttpLoadService{}
	if err := r.Get(ctx, req.NamespacedName, loadservice); err != nil {
		return ctrl.Result{}, err
	}

	if err := usecase.Load(*loadservice.DeepCopy()); err != nil {
		r.Log.Error(err, "load use case error")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, err
	}

	return ctrl.Result{}, nil
}

func (r *HttpLoadServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&loadv1alpha1.HttpLoadService{}).
		Complete(r)
}
