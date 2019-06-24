package stop

import (
	"time"

	"github.com/brad-jones/goasync/task"
)

type Stopable interface {
	Stop()
}

func ToStopables(awaitables ...interface{}) []Stopable {
	out := []Stopable{}
	for _, v := range awaitables {
		out = append(out, v.(Stopable))
	}
	return out
}

func All(awaitables ...Stopable) {
	for _, awaitable := range awaitables {
		awaitable.Stop()
	}
}

func AllAsync(awaitables ...Stopable) *task.Task {
	return task.New(func(t *task.TaskInternal) {
		All(awaitables...)
		t.Resolve(struct{}{})
	})
}

type StopableWithTimeout interface {
	StopWithTimeout(timeout time.Duration) error
}

func ToStopablesWithTimeout(awaitables ...interface{}) []StopableWithTimeout {
	out := []StopableWithTimeout{}
	for _, v := range awaitables {
		out = append(out, v.(StopableWithTimeout))
	}
	return out
}

func AllWithTimeout(timeout time.Duration, awaitables ...StopableWithTimeout) {
	for _, awaitable := range awaitables {
		awaitable.StopWithTimeout(timeout)
	}
}

func AllWithTimeoutAsync(timeout time.Duration, awaitables ...StopableWithTimeout) *task.Task {
	return task.New(func(t *task.TaskInternal) {
		AllWithTimeout(timeout, awaitables...)
		t.Resolve(struct{}{})
	})
}
