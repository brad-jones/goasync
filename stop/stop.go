// Package stop is a means of stopping many tasks in bulk.
package stop

import (
	"time"

	"github.com/brad-jones/goasync/v2/task"
)

// All will loop through all provided objects and call their Stop method.
func All(stopables ...*task.Task) {
	for _, stopable := range stopables {
		stopable.Stop()
	}
}

// AllAsync does exactly the same thing as All but does so asynchronously.
func AllAsync(stopables ...*task.Task) *task.Task {
	return task.New(func(t *task.Internal) {
		All(stopables...)
	})
}

// AllWithTimeout all provided objects and call their StopWithTimeout
// method with the given timeout value.
func AllWithTimeout(timeout time.Duration, stopables ...*task.Task) {
	for _, stopable := range stopables {
		stopable.StopWithTimeout(timeout)
	}
}

// AllWithTimeoutAsync does exactly the same thing as AllWithTimeout but does so asynchronously.
func AllWithTimeoutAsync(timeout time.Duration, stopables ...*task.Task) *task.Task {
	return task.New(func(t *task.Internal) {
		AllWithTimeout(timeout, stopables...)
	})
}
