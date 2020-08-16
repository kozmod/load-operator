package executor

import (
	"context"
	"time"
)

const defaultDuration = 1 * time.Second

type UtcScheduleExecutor struct {
	start     *time.Time
	duration  time.Duration
	executeFn func()
}

type Option func(*UtcScheduleExecutor)

func WithUtcStart(start time.Time) Option {
	return func(executor *UtcScheduleExecutor) {
		utc := asUtc(start)
		executor.start = &utc
	}
}

func NewScheduleExecutor(duration time.Duration, execute func(), opts ...Option) ScheduleExecutor {
	executor := &UtcScheduleExecutor{duration: duration, executeFn: execute}
	for _, opt := range opts {
		opt(executor)
	}
	return executor
}

func (s *UtcScheduleExecutor) Schedule(ctx context.Context) {
	currentTime := time.Now().UTC()
	duration := betweenOrDefault(&currentTime, s.start, defaultDuration)
	duration = defaultIfNotPositive(duration, defaultDuration)
	timer := time.NewTimer(duration)
	for {
		select {
		case <-timer.C:
			timer = time.NewTimer(s.duration)
			s.executeFn()
		case <-ctx.Done():
			timer.Stop()
			return
		}
	}
}

func betweenOrDefault(start *time.Time, end *time.Time, def time.Duration) time.Duration {
	if start == nil || end == nil {
		return def
	}
	return end.Sub(*start)
}

func defaultIfNotPositive(duration time.Duration, def time.Duration) time.Duration {
	if duration <= 0 {
		duration = def
	}
	return duration
}

func asUtc(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
}
