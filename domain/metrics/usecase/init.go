package usecase

import (
	"github.com/go-logr/logr"
	"github.com/kozmod/load-operator/domain/metrics/usecase/internal"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var uc *ScheduleUseCase

func Init(conf *rest.Config, cl client.Client, l logr.Logger) *ScheduleUseCase {
	if uc == nil {
		return NewScheduleUseCase(
			internal.NewPrint(
				*internal.NewMetrics(conf, cl, l),
				l,
			), l)
	}
	return uc
}
