package usecase

import (
	"context"

	"github.com/kozmod/load-operator/apis/load/v1alpha1"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type loadUC interface {
	Load(v1alpha1.HttpLoadService) error
	Metrics() (vegeta.Metrics, error)
}

type scheduleUC interface {
	Schedule(context.Context, v1alpha1.MetricsService) error
}
