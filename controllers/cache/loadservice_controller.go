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
	"fmt"
	"github.com/go-logr/logr"
	cachev1 "github.com/kozmod/load-operator/apis/cache/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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
	_ = r.Log.WithValues("loadservice", req.NamespacedName)

	// your logic here
	loadService := &cachev1.LoadService{}
	if err := r.Get(ctx, req.NamespacedName, loadService); err != nil {
		return ctrl.Result{}, err
	}
	r.Log.Info("loadService:\n", "l", loadService)

	deployment := &appsv1.Deployment{}
	if err := r.Get(ctx, client.ObjectKey{
		Namespace: loadService.Spec.Namespace,
		Name:      loadService.Spec.DeploymentName,
	}, deployment); err != nil {
		r.Log.Error(err, "Failed to list deployment.")
		return ctrl.Result{RequeueAfter: time.Second * 5}, err
	}

	const appLabel = "app"
	pods := &corev1.PodList{}
	opts := []client.ListOption{
		client.InNamespace(loadService.Spec.Namespace),
		client.MatchingLabels{appLabel: deployment.Spec.Selector.MatchLabels[appLabel]}}
	if err := r.List(ctx, pods, opts...); err != nil {
		r.Log.Error(err, "Failed to get pods.")
		return ctrl.Result{RequeueAfter: time.Second * 5}, err
	}

	config := rest.AnonymousClientConfig(&rest.Config{Host: "127.0.0.1:61866"}) //todo local test
	//config := rest.InClusterConfig() //todo in cluster
	//if err != nil {
	//	fmt.Println("rest error")
	//	panic(err.Error())
	//}
	clientset, err := metricsv.NewForConfig(config)
	if err != nil {
		fmt.Println("metricsv.NewForConfig(config) error")
		panic(err.Error())
	}
	for i := 0; i < len(pods.Items); i++ {
		podMetrics, err := clientset.MetricsV1beta1().PodMetricses(loadService.Spec.Namespace).
			Get(context.TODO(), pods.Items[i].Name, metav1.GetOptions{})
		if err != nil {
			r.Log.Error(err, "get metrics error")
		}
		r.Log.Info("pod metrics", "metric "+string(rune(i)), podMetrics)
	}

	return ctrl.Result{}, nil
}

func (r *LoadServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1.LoadService{}).
		Complete(r)
}
