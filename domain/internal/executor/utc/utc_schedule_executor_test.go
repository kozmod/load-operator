package utc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBetweenOrDefault_DefaultWhenBothNil(t *testing.T) {
	expected := 1 * time.Nanosecond
	require.Equal(t, betweenOrDefault(nil, nil, expected), expected)
}

func TestBetweenOrDefault_DefaultWhenStartNil(t *testing.T) {
	expected := 1 * time.Nanosecond
	current := time.Now()
	require.Equal(t, betweenOrDefault(nil, &current, expected), expected)
}

func TestBetweenOrDefault_DefaultWhenEndNil(t *testing.T) {
	expected := 1 * time.Nanosecond
	current := time.Now()
	require.Equal(t, betweenOrDefault(&current, nil, expected), expected)
}

func TestBetweenOrDefault_Result(t *testing.T) {
	expected := 0 * time.Nanosecond
	current := time.Now()
	require.Equal(t, betweenOrDefault(&current, &current, expected), expected)
}

func TestBetweenOrDefault_PositiveAndNegativeResult(t *testing.T) {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fail()
	}
	start := time.Date(2021, time.Month(1), 23, 12, 30, 05, 5, location)
	end := time.Date(2020, time.Month(1), 23, 12, 30, 05, 5, location)
	require.Equal(t, betweenOrDefault(&start, &end, 0*time.Second), -8784*time.Hour)

	start = time.Date(2020, time.Month(1), 23, 12, 30, 05, 5, location)
	end = time.Date(2021, time.Month(1), 23, 12, 30, 05, 5, location)
	require.Equal(t, betweenOrDefault(&start, &end, 0*time.Second), 8784*time.Hour)
}

func TestDefaultIfNotPositive_DefaultWhenNegative(t *testing.T) {
	val := -1 * time.Second
	def := 10 * time.Hour
	require.Equal(t, defaultIfNotPositive(val, def), def)
}

func TestDefaultIfNotPositive_DefaultWhenZero(t *testing.T) {
	val := 0 * time.Second
	def := 10 * time.Second
	require.Equal(t, defaultIfNotPositive(val, def), def)
}

func TestDefaultIfNotPositive_Val(t *testing.T) {
	val := 1 * time.Second
	def := 10 * time.Second
	require.Equal(t, defaultIfNotPositive(val, def), val)
}

func TestExecutor_ShouldStopByCancel(t *testing.T) {
	start := time.Now().Add(1 * time.Second)
	duration := 1 * time.Second
	ctx, cancel := context.WithCancel(context.TODO())
	es := &ScheduleExecutor{start: &start, duration: duration,
		executeFn: func() {
			t.Fatal()
		}}
	go cancel()
	es.Schedule(ctx)
	<-ctx.Done()
}

func TestExecutor_ShouldStopByTimeout(t *testing.T) {
	start := time.Now().Add(1 * time.Second)
	duration := 1 * time.Second
	res := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Microsecond)
	es := &ScheduleExecutor{start: &start, duration: duration,
		executeFn: func() {
			res <- struct{}{}
		}}
	go es.Schedule(ctx)
	select {
	case <-res:
		t.Fail()
	case <-ctx.Done():
		break
	}
	cancel()
}

func TestOption_WithUtcStart_CastCurrentToUtc(t *testing.T) {
	location, _ := time.LoadLocation("America/New_York")
	locationDate := time.Date(2020, 1, 1, 8, 1, 1, 1, location)

	es := &ScheduleExecutor{}
	option := WithUtcStart(locationDate)
	option(es)
	require.NotNil(t, es.duration)
	require.NotNil(t, es.start)
	require.Equal(t, locationDate.Year(), es.start.Year())
	require.Equal(t, locationDate.Month(), es.start.Month())
	require.Equal(t, locationDate.Day(), es.start.Day())
	require.Equal(t, locationDate.Hour(), es.start.Hour())
	require.Equal(t, locationDate.Minute(), es.start.Minute())
	require.Equal(t, locationDate.Second(), es.start.Second())
	require.Equal(t, locationDate.Nanosecond(), es.start.Nanosecond())
	require.Equal(t, time.UTC, es.start.Location())
}

func TestOption_WithUtcStart(t *testing.T) {
	current := time.Now()
	es := &ScheduleExecutor{}
	option := WithUtcStart(current)
	option(es)
	require.NotNil(t, es.duration)
	require.NotNil(t, es.start)
	require.Equal(t, current.Year(), es.start.Year())
	require.Equal(t, current.Month(), es.start.Month())
	require.Equal(t, current.Day(), es.start.Day())
	require.Equal(t, current.Hour(), es.start.Hour())
	require.Equal(t, current.Minute(), es.start.Minute())
	require.Equal(t, current.Second(), es.start.Second())
	require.Equal(t, current.Nanosecond(), es.start.Nanosecond())
	require.Equal(t, time.UTC, es.start.Location())
}

func TestAsUtc(t *testing.T) {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fail()
	}
	src := time.Date(2021, time.Month(1), 23, 12, 30, 05, 5, location)
	res := asUtc(src)
	require.Equal(t, src.Year(), res.Year())
	require.Equal(t, src.Month(), res.Month())
	require.Equal(t, src.Day(), res.Day())
	require.Equal(t, src.Hour(), res.Hour())
	require.Equal(t, src.Minute(), res.Minute())
	require.Equal(t, src.Second(), res.Second())
	require.Equal(t, src.Nanosecond(), res.Nanosecond())
	require.Equal(t, time.UTC, res.Location())
}
