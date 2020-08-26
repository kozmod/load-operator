package usecase

import (
	"context"
	"github.com/go-logr/logr"
	cachev1 "github.com/kozmod/load-operator/apis/cache/v1"
	"github.com/kozmod/load-operator/domain/metrics/usecase/internal/get"
	"github.com/kozmod/load-operator/domain/metrics/usecase/internal/load"
	print2 "github.com/kozmod/load-operator/domain/metrics/usecase/internal/print"
	"github.com/kozmod/load-operator/domain/metrics/usecase/internal/schedule"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var muc *MetricsUseCaseInteractor

type MetricsUseCaseInteractor struct {
	suc scheduleExecutor
	luc loader
}

func (uc *MetricsUseCaseInteractor) Apply(ctx context.Context, loadService cachev1.LoadService) error {
	err := uc.suc.Schedule(ctx, loadService)
	if err != nil {
		return err
	}
	err = uc.luc.Load(loadService)
	if err != nil {
		return err
	}
	return nil
}

func TryInit(conf *rest.Config, cl client.Client, l logr.Logger) *MetricsUseCaseInteractor {
	if muc == nil {
		luc := load.NewLoader()
		muc = &MetricsUseCaseInteractor{
			luc: luc,
			suc: schedule.NewScheduleUseCase(
				print2.NewPrint(
					get.NewMetrics(conf, cl, l),
					luc,
					l,
				), l),
		}
	}
	return muc
}
