package task

//go:generate genny -in=generic.genny -out=typed.go gen "Awaitable=BUILTINS"
//go:generate genny -in=generic_slice.genny -out=typed_slice.go gen "Awaitable=BUILTINS"

import (
	"time"

	"github.com/brad-jones/goerr"
	"github.com/go-errors/errors"
)

// Task represents an asynchronous operation, create new instances with New.
//
// Not type safe, use with care.
type Task struct {
	Resolver <-chan interface{}
	Rejector <-chan error
	Stopper  *chan struct{}
	closed   bool
	done     chan struct{}
	value    interface{}
	err      error
}

func (t *Task) Stop() {
	if !t.closed {
		defer func() { recover() }()
		close(*t.Stopper)
		t.closed = true
	}
	<-t.done
}

func (t *Task) StopWithTimeout(timeout time.Duration) error {
	if !t.closed {
		defer func() { recover() }()
		close(*t.Stopper)
		t.closed = true
	}
	select {
	case <-t.done:
		return nil
	case <-time.After(timeout):
		return errors.New(&ErrStoppingTaskTimeout{})
	}
}

func (t *Task) Result() (interface{}, error) {
	<-t.done
	return t.value, t.err
}

func (t *Task) MustResult() interface{} {
	v, e := t.Result()
	goerr.Check(e)
	return v
}

// TaskInternal is passed into a new task function.
//
// Not type safe, use with care.
type TaskInternal struct {
	Resolver chan<- interface{}
	Rejector chan<- error
	Stopper  *chan struct{}
}

func (ti *TaskInternal) Resolve(v interface{}) {
	ti.Resolver <- v
}

func (ti *TaskInternal) Reject(e error) {
	ti.Rejector <- errors.Wrap(e, 0)
}

func (ti *TaskInternal) ShouldStop() bool {
	select {
	case <-*ti.Stopper:
		return true
	default:
		return false
	}
}

// New creates new instances of Task.
//
// Not type safe, use with care.
func New(fn func(t *TaskInternal)) *Task {
	// Spin up some channels
	done := make(chan struct{}, 1)
	stopper := make(chan struct{}, 1)
	tResolver := make(chan interface{}, 1)
	tRejector := make(chan error, 1)
	tiResolver := make(chan interface{}, 1)
	tiRejector := make(chan error, 1)

	// Pass those channels into our tasks
	t := &Task{
		Resolver: tResolver,
		Rejector: tRejector,
		Stopper:  &stopper,
		done:     done,
	}
	ti := &TaskInternal{
		Resolver: tiResolver,
		Rejector: tiRejector,
		Stopper:  &stopper,
	}

	// Execute the task asynchronously
	go func() {
		// Regardless of what the function does we know that it is done
		defer close(done)

		// Catch any panics and reject them
		defer goerr.Handle(func(e error) {
			t.err = e
			tRejector <- t.err
		})

		// Execute the task
		fn(ti)

		// Read the result in a non blocking manner
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
//
// Not type safe, use with care.
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
		done:     done,
		closed:   true,
	}
}

// Rejected returns a pre-rejected task
//
// Not type safe, use with care.
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
		done:     done,
		closed:   true,
	}
}
