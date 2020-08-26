package print

import (
	"context"

	vegeta "github.com/tsenart/vegeta/v12/lib"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"

	loadv1alpha1 "github.com/kozmod/load-operator/apis/load/v1alpha1"
)

type loader interface {
	Metrics() (vegeta.Metrics, error)
}

type metrics interface {
	Get(context.Context, loadv1alpha1.MetricsService) (map[*corev1.Pod]*v1beta1.PodMetrics, error)
}
