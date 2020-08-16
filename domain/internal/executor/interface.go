package executor

import "context"

type ScheduleExecutor interface {
	Schedule(ctx context.Context)
}
