package print

import (
	"context"

	"github.com/go-logr/logr"
	loadv1alpha1 "github.com/kozmod/load-operator/apis/load/v1alpha1"
)

type PrintlnMetrics struct {
	k8sMetrics    metrics
	loaderMetrics loader
	log           logr.Logger
}

func New(m metrics, ldr loader, l logr.Logger) *PrintlnMetrics {
	return &PrintlnMetrics{
		k8sMetrics:    m,
		loaderMetrics: ldr,
		log:           l,
	}
}

func (uc *PrintlnMetrics) Apply(ctx context.Context, m loadv1alpha1.MetricsService) error {
	metrics, err := uc.k8sMetrics.Get(ctx, m)
	if err != nil {
		return err
	}
	lm, err := uc.loaderMetrics.Metrics()
	if err != nil {
		return err
	}
	//uc.log.Info("loader metrics\n", "loader", lm)
	uc.log.Info("loader metrics\n", "loader", lm)
	for k, v := range metrics {
		uc.log.Info("metrics\n", "pod", *k)
		uc.log.Info("metrics\n", "metrics", *v)
	}
	return nil
}
