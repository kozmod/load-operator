package print

import (
	"context"
	"github.com/go-logr/logr"
	v1 "github.com/kozmod/load-operator/apis/cache/v1"
)

type PrintlnMetrics struct {
	k8sMetrics    metrics
	loaderMetrics loader
	log           logr.Logger
}

func NewPrint(m metrics, ldr loader, l logr.Logger) *PrintlnMetrics {
	return &PrintlnMetrics{
		k8sMetrics:    m,
		loaderMetrics: ldr,
		log:           l,
	}
}

func (uc *PrintlnMetrics) Apply(ctx context.Context, loadService v1.LoadService) error {
	metrics, err := uc.k8sMetrics.Get(ctx, loadService)
	if err != nil {
		return err
	}
	uc.log.Info("loader metrics\n", "loader", uc.loaderMetrics.Metrics())
	for k, v := range metrics {
		uc.log.Info("metrics\n", "pod", *k)
		uc.log.Info("metrics\n", "metrics", *v)
	}
	return nil
}
