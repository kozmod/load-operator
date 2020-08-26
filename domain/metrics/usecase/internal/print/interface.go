package print

import (
	"context"
	cachev1 "github.com/kozmod/load-operator/apis/cache/v1"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type loader interface {
	Metrics() vegeta.Metrics
}

type metrics interface {
	Get(ctx context.Context, loadService cachev1.LoadService) (map[*corev1.Pod]*v1beta1.PodMetrics, error)
}
