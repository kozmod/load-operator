package internal

import (
	"context"
	"github.com/go-logr/logr"
	v1 "github.com/kozmod/load-operator/apis/cache/v1"
)

type PrintlnMetrics struct {
	metrics Metrics
	log     logr.Logger
}

func NewPrint(m Metrics, l logr.Logger) *PrintlnMetrics {
	return &PrintlnMetrics{
		metrics: m,
		log:     l,
	}
}

func (uc *PrintlnMetrics) Apply(ctx context.Context, loadService v1.LoadService) error {
	metrics, err := uc.metrics.Get(ctx, loadService)
	if err != nil {
		return err
	}
	for k, v := range metrics {
		uc.log.Info("metrics\n", "pod", *k)
		uc.log.Info("metrics\n", "metrics", *v)
	}
	return nil
}
