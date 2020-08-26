package schedule

import (
	"context"
	"github.com/kozmod/load-operator/apis/cache/v1"
)

type useCase interface {
	Apply(ctx context.Context, loadService v1.LoadService) error
}
