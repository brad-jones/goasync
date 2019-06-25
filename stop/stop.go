/*
Package stop is a means of stopping many tasks in bulk.

Find main reference documentation at https://godoc.org/github.com/brad-jones/goasync
*/
package stop

import (
	"time"

	"github.com/brad-jones/goasync/task"
)

type Stopable interface {
	Stop()
}

func All(stopables ...Stopable) {
	for _, stopable := range stopables {
		stopable.Stop()
	}
}

func AllAsync(stopables ...Stopable) *task.Task {
	return task.New(func(t *task.Internal) {
		All(stopables...)
	})
}

type StopableWithTimeout interface {
	StopWithTimeout(timeout time.Duration) error
}

func AllWithTimeout(timeout time.Duration, stopables ...StopableWithTimeout) {
	for _, stopable := range stopables {
		stopable.StopWithTimeout(timeout)
	}
}

func AllWithTimeoutAsync(timeout time.Duration, stopables ...StopableWithTimeout) *task.Task {
	return task.New(func(t *task.Internal) {
		AllWithTimeout(timeout, stopables...)
	})
}
