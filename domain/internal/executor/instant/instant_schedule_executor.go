package instant

import (
	"context"
	"time"
)

type ScheduleExecutor struct {
	immediately bool
	duration    time.Duration
	executeFn   func()
}

type Option func(*ScheduleExecutor)

func Immediately() Option {
	return func(executor *ScheduleExecutor) {
		executor.immediately = true
	}
}

func NewScheduleExecutor(duration time.Duration, execute func(), opts ...Option) *ScheduleExecutor {
	e := &ScheduleExecutor{duration: duration, executeFn: execute}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (s *ScheduleExecutor) Schedule(ctx context.Context) {
	if s.immediately {
		go s.executeFn()
	}
	ticker := time.NewTicker(s.duration)
	for {
		select {
		case <-ticker.C:
			go s.executeFn()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
