package schedule

import (
	"context"
	"github.com/go-logr/logr"
	cachev1 "github.com/kozmod/load-operator/apis/cache/v1"
	"github.com/kozmod/load-operator/domain/internal/executor/instant"
)

type UseCase struct {
	useCase    useCase
	cancelFunc context.CancelFunc
	log        logr.Logger
}

func NewScheduleUseCase(uc useCase, l logr.Logger) *UseCase {
	return &UseCase{
		useCase: uc,
		log:     l,
	}
}

func (s *UseCase) Schedule(ctx context.Context, loadService cachev1.LoadService) error {
	ctx, cancel := context.WithCancel(ctx)
	ts := instant.NewScheduleExecutor(loadService.Spec.Metrics.Duration.Duration, func() {
		err := s.useCase.Apply(ctx, *loadService.DeepCopy())
		if err != nil {
			s.log.Error(err, "schedule execute error")
		}
	})
	s.Stop()
	s.cancelFunc = cancel
	go ts.Schedule(ctx)
	return nil
}

func (s *UseCase) Stop() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}
