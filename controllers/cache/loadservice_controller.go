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
	cachev1 "github.com/kozmod/load-operator/apis/cache/v1"
	"github.com/kozmod/load-operator/domain/metrics/usecase"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const appLabel = "app"

// LoadServiceReconciler reconciles a LoadService object
type LoadServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cache.load-operator.com,resources=loadservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.load-operator.com,resources=loadservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

func (r *LoadServiceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	l := r.Log.WithValues("loadservice", req.NamespacedName)

	// your logic here
	loadService := &cachev1.LoadService{}
	if err := r.Get(ctx, req.NamespacedName, loadService); err != nil {
		return ctrl.Result{}, err
	}
	config := rest.AnonymousClientConfig(&rest.Config{Host: "127.0.0.1:64931"}) //todo local test
	//config := rest.InClusterConfig() //todo in cluster
	//if err != nil {
	//	fmt.Println("rest error")
	//	panic(err.Error())
	//}
	uc := usecase.TryInit(config, r.Client, l)

	if err := uc.Apply(ctx, *loadService.DeepCopy()); err != nil {
		r.Log.Error(err, "usecase error")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, err
	}

	return ctrl.Result{}, nil
}

func (r *LoadServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1.LoadService{}).
		Complete(r)
}
