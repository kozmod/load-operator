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
	"github.com/go-logr/logr"
	"github.com/kozmod/load-operator/domain/metrics/usecase"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	loadv1alpha1 "github.com/kozmod/load-operator/apis/load/v1alpha1"
)

// MetricsServiceReconciler reconciles a MetricsService object
type MetricsServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=load.load-operator.com,resources=metricsservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=load.load-operator.com,resources=metricsservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

func (r *MetricsServiceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	l := r.Log.WithValues("loadservice", req.NamespacedName)

	// your logic here
	metrics := &loadv1alpha1.MetricsService{}
	if err := r.Get(ctx, req.NamespacedName, metrics); err != nil {
		return ctrl.Result{}, err
	}
	usecase.InitScheduleUseCase(r.Client, l)
	if err := usecase.Schedule(ctx, *metrics.DeepCopy()); err != nil {
		r.Log.Error(err, "use case error")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, err
	}

	return ctrl.Result{}, nil
}

func (r *MetricsServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&loadv1alpha1.MetricsService{}).
		Complete(r)
}
