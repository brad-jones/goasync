// Package task is a asynchronous utility inspired by JS Promises & C# Tasks.
package task

import (
	"context"
	"time"

	"github.com/brad-jones/goerr/v2"
)

// Task represents an asynchronous operation, create new instances with New.
type Task struct {
	// Every task has a resolver channel that ultimately represents a
	// value to be returned sometime in the future. Keep in mind that
	// once a value is read from this channel it can not be read from again.
	Resolver <-chan interface{}

	// Every task has a rejector channel that represents a possible error
	// to be returned sometime in the future. Keep in mind that once a value
	// is read from this channel it can not be read from again.
	Rejector <-chan error

	// Every task has a stopper channel that allows the task to be
	// stopped cooperatively from the outside. Simply close this
	// channel to stop the task.
	Stopper *chan struct{}

	// Used to track when the task has actually finished regardless of what
	// has or hasn't been resolved/rejected. Tasks can just perform some action
	// they don't necessarily have to resolve something.
	Done *chan struct{}

	// We keep a copy of the value for use with Result()
	value interface{}

	// We keep a copy of the error for use with Result()
	err error
}

// Stop the task cooperatively, this will block until the task has returned.
func (t *Task) Stop() {
	defer func() { recover() }()
	close(*t.Stopper)
	<-*t.Done
}

// StopWithTimeout will stop the task cooperatively but return an error if
// a timeout is reached. Use this to ensure your application does not
// hang indefinitely.
func (t *Task) StopWithTimeout(timeout time.Duration) error {
	defer func() { recover() }()
	close(*t.Stopper)

	select {
	case <-*t.Done:
		return nil
	case <-time.After(timeout):
		return goerr.Wrap(&ErrStoppingTaskTimeout{})
	}
}

// ErrStoppingTaskTimeout is returned by StopWithTimeout when the given duration
// has passed and the task has still not stopped.
type ErrStoppingTaskTimeout struct {
}

func (e *ErrStoppingTaskTimeout) Error() string {
	return "task: stopping task took too long to stop"
}

// Result waits for the task to complete and then returns any resolved
// (or rejected) values. This can be called many times over and the same
// values will be returned.
func (t *Task) Result() (interface{}, error) {
	<-*t.Done
	return t.value, t.err
}

// MustResult does the same as Result() but panics if an error was rejected.
func (t *Task) MustResult() interface{} {
	v, e := t.Result()
	goerr.Check(e)
	return v
}

// ResultWithTimeout takes 2 duration values, the first is the amount of time
// we will wait for the task to complete, if that time passes we will then call
// `StopWithTimeout` which will wait for the second duration for the given task
// to cooperatively stop.
func (t *Task) ResultWithTimeout(runtime, stoptime time.Duration) (interface{}, error) {
	select {
	case <-*t.Done:
		return t.value, t.err
	case <-time.After(runtime):
		if err := t.StopWithTimeout(stoptime); err != nil {
			return nil, goerr.Wrap(err, "result took too long to be returned and failed to stop in a timely manner")
		}
		return nil, goerr.Wrap(t.err, "result took too long to be returned")
	}
}

// MustResultWithTimeout does the same as ResultWithTimeout() but panics if an error was encountered.
func (t *Task) MustResultWithTimeout(runtime, stoptime time.Duration) interface{} {
	v, e := t.ResultWithTimeout(runtime, stoptime)
	goerr.Check(e)
	return v
}

// Internal is used by the task implementor.
type Internal struct {
	// Every task has a resolver channel that ultimately represents a
	// value to be returned sometime in the future. As a task implementor
	// you should only ever send a single value to this channel.
	Resolver chan<- interface{}

	// Every task has a rejector channel that represents a possible error
	// to be returned sometime in the future. As a task implementor
	// you should only ever send a single value to this channel.
	Rejector chan<- error

	// Every task has a stopper channel that allows the task to be
	// stopped cooperatively from the outside. If this channel is closed
	// you should stop your task.
	Stopper *chan struct{}

	// Used internally to track when the task has actually finished
	// regardless of what has or hasn't been resolved/rejected.
	done *chan struct{}
}

// Resolve is a simple function that sends the provided value to the resolver channel.
func (i *Internal) Resolve(v interface{}) {
	i.Resolver <- v
}

// Reject is a simple function that sends the provided error to the rejector channel.
func (i *Internal) Reject(err interface{}, messages ...string) {
	i.Rejector <- goerr.Trace(1, err, messages...)
}

// ShouldStop is a non blocking method that informs your task if it should stop.
func (i *Internal) ShouldStop() bool {
	select {
	case <-*i.Stopper:
		return true
	default:
		return false
	}
}

// CancelableCtx returns a context object that will be canceled if this task is
// told to stop, this is useful for integrating with more traditional go code.
func (i *Internal) CancelableCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-*i.done:
		case <-*i.Stopper:
			cancel()
		}
	}()
	return ctx
}

// New creates new instances of Task.
func New(fn func(t *Internal)) *Task {
	// Spin up some channels
	done := make(chan struct{}, 1)
	stopper := make(chan struct{}, 1)
	tResolver := make(chan interface{}, 1)
	tRejector := make(chan error, 1)
	tiResolver := make(chan interface{}, 1)
	tiRejector := make(chan error, 1)

	// Create our task object
	t := &Task{
		Resolver: tResolver,
		Rejector: tRejector,
		Stopper:  &stopper,
		Done:     &done,
	}

	// Execute the task asynchronously
	go func() {
		// Regardless of what the function does we know that it is done
		defer close(Done)

		// Catch any panics and reject them
		defer goerr.Handle(func(e error) {
			t.err = goerr.Trace(3, e)
			tRejector <- t.err
		})

		// Execute the task
		fn(&Internal{
			Resolver: tiResolver,
			Rejector: tiRejector,
			Stopper:  t.Stopper,
			done:     &done,
		})

		// Read the result in a non blocking manner. Keep in mind not every task
		// will actually resolve or reject anything, the simple fact that it is
		// done could be enough.
		select {
		case v := <-tiResolver:
			t.value = v
			tResolver <- t.value
		case e := <-tiRejector:
			t.err = e
			tRejector <- t.err
		default:
		}
	}()

	// Return the task object
	return t
}

// Resolved returns a pre-resolved task
func Resolved(v interface{}) *Task {
	done := make(chan struct{}, 1)
	close(done)
	resolver := make(chan interface{}, 1)
	resolver <- v
	rejector := make(chan error, 1)
	return &Task{
		Resolver: resolver,
		Rejector: rejector,
		Stopper:  &done,
		Done:     &done,
	}
}

// Rejected returns a pre-rejected task
func Rejected(e error) *Task {
	done := make(chan struct{}, 1)
	close(done)
	resolver := make(chan interface{}, 1)
	rejector := make(chan error, 1)
	rejector <- e
	return &Task{
		Resolver: resolver,
		Rejector: rejector,
		Stopper:  &done,
		Done:     &done,
	}
}
