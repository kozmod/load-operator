package schedule

import (
	"context"
	loadv1alpha1 "github.com/kozmod/load-operator/apis/load/v1alpha1"
)

type useCase interface {
	Apply(context.Context, loadv1alpha1.MetricsService) error
}
