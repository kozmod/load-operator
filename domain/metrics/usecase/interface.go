package usecase

import (
	"context"
	cachev1 "github.com/kozmod/load-operator/apis/cache/v1"
)

type loader interface {
	Load(loadService cachev1.LoadService) error
}

type scheduleExecutor interface {
	Schedule(context.Context, cachev1.LoadService) error
}
