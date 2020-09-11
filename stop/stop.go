// Package stop is a means of stopping many tasks in bulk.
package stop

import (
	"time"

	"github.com/brad-jones/goasync/v2/task"
)

// Stopable is any object that has a Stop method,
// the functions in this package then call that method for you.
type Stopable interface {
	Stop()
}

// All will loop through all provided objects and call their Stop method.
func All(stopables ...Stopable) {
	for _, stopable := range stopables {
		stopable.Stop()
	}
}

// AllAsync does exactly the same thing as All but does so asynchronously.
func AllAsync(stopables ...Stopable) *task.Task {
	return task.New(func(t *task.Internal) {
		All(stopables...)
	})
}

// StopableWithTimeout is any object that has a StopWithTimeout method,
// the functions in this package then call that method for you.
type StopableWithTimeout interface {
	StopWithTimeout(timeout time.Duration) error
}

// AllWithTimeout all provided objects and call their StopWithTimeout
// method with the given timeout value.
func AllWithTimeout(timeout time.Duration, stopables ...StopableWithTimeout) {
	for _, stopable := range stopables {
		stopable.StopWithTimeout(timeout)
	}
}

// AllWithTimeoutAsync does exactly the same thing as AllWithTimeout but does so asynchronously.
func AllWithTimeoutAsync(timeout time.Duration, stopables ...StopableWithTimeout) *task.Task {
	return task.New(func(t *task.Internal) {
		AllWithTimeout(timeout, stopables...)
	})
}
