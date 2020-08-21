package instant

import (
	"context"
	"testing"
	"time"
)

func TestExecutor_ShouldStopByTimeout(t *testing.T) {
	duration := 1 * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Nanosecond)
	res := make(chan struct{})
	es := &ScheduleExecutor{duration: duration,
		executeFn: func() {
			res <- struct{}{}
		}}
	go es.Schedule(ctx)
	select {
	case <-res:
		t.Fail()
	case <-ctx.Done():
		cancel()
	}
}
