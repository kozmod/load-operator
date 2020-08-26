package get

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	loadv1alpha1 "github.com/kozmod/load-operator/apis/load/v1alpha1"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const appLabel = "app"

type UseCase struct {
	conf   *rest.Config
	client client.Client
	log    logr.Logger
}

func New(conf *rest.Config, cl client.Client, l logr.Logger) *UseCase {
	return &UseCase{
		conf:   conf,
		client: cl,
		log:    l,
	}
}

func (s *UseCase) Get(ctx context.Context, ms loadv1alpha1.MetricsService) (map[*corev1.Pod]*v1beta1.PodMetrics, error) {
	namespace := ms.Spec.Namespace
	deploymentName := ms.Spec.DeploymentName
	pods, err := s.getPodList(ctx, namespace, deploymentName)
	if err != nil {
		return nil, err
	}
	metrics, err := s.getMetrics(ctx, pods, namespace)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (s *UseCase) getPodList(ctx context.Context, namespace, deploymentName string) (*corev1.PodList, error) {
	deployment := &appsv1.Deployment{}
	if err := s.client.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      deploymentName,
	}, deployment); err != nil {
		return nil, errors.WithMessage(err, "get deployment error")
	}

	pods := &corev1.PodList{}
	opts := []client.ListOption{
		client.InNamespace(namespace),
		client.MatchingLabels{appLabel: deployment.Spec.Selector.MatchLabels[appLabel]}}
	if err := s.client.List(ctx, pods, opts...); err != nil {
		return nil, errors.WithMessage(err, "get pod list error")
	}
	return pods, nil
}

func (s *UseCase) getMetrics(ctx context.Context, pods *corev1.PodList, namespace string) (map[*corev1.Pod]*v1beta1.PodMetrics, error) {
	clientset, err := metricsv.NewForConfig(s.conf)
	if err != nil {
		return nil, errors.WithMessage(err, "create clientset error")
	}
	metricsmap := make(map[*corev1.Pod]*v1beta1.PodMetrics, 0)
	for i := 0; i < len(pods.Items); i++ {
		pod := &pods.Items[i]
		podMetrics, err := clientset.MetricsV1beta1().
			PodMetricses(namespace).
			Get(ctx, pod.Name, metav1.GetOptions{})
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("get pod metrics error: %s", pod.ObjectMeta.Name))
		}
		metricsmap[pod] = podMetrics
	}
	return metricsmap, nil
}
