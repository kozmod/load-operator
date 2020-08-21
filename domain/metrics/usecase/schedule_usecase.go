package usecase

import (
	"context"
	"github.com/go-logr/logr"
	cachev1 "github.com/kozmod/load-operator/apis/cache/v1"
	"github.com/kozmod/load-operator/domain/internal/executor/instant"
)

type ScheduleUseCase struct {
	useCase    useCase
	cancelFunc context.CancelFunc
	log        logr.Logger
}

func NewScheduleUseCase(uc useCase, l logr.Logger) *ScheduleUseCase {
	return &ScheduleUseCase{
		useCase: uc,
		log:     l,
	}
}

func (s *ScheduleUseCase) Schedule(ctx context.Context, loadService cachev1.LoadService) error {
	ctx, cancel := context.WithCancel(ctx)
	ts := instant.NewScheduleExecutor(loadService.Spec.Delay.Duration, func() {
		err := s.useCase.Apply(ctx, *loadService.DeepCopy())
		if err != nil {
			s.log.Error(err, "schedule execute error")
		}
	}, instant.Immediately())
	s.Stop()
	s.cancelFunc = cancel
	go ts.Schedule(ctx)
	return nil
}

func (s *ScheduleUseCase) Stop() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}
