/*
Package await contains await helper functions for use with tasks.

Find main reference documentation at https://godoc.org/github.com/brad-jones/goasync
*/
package await

import (
	"time"

	"github.com/brad-jones/goasync/stop"
	"github.com/go-errors/errors"
)

// Awaitable refers to any type that has a Result() method
type Awaitable interface {
	Result() (interface{}, error)
}

// All will wait for every given task to emit a result, the results (& errors)
// will be returned in a slice ordered the same as the input.
func All(awaitables ...Awaitable) ([]interface{}, error) {
	awaited := []interface{}{}
	awaitedErrors := []error{}

	for _, awaitable := range awaitables {
		v, e := awaitable.Result()
		if e != nil {
			awaitedErrors = append(awaitedErrors, errors.Wrap(e, 0))
		}
		awaited = append(awaited, v)
	}

	if len(awaitedErrors) > 0 {
		return nil, errors.New(&ErrTaskFailed{
			Errors: awaitedErrors,
		})
	}

	return awaited, nil
}

// AllOrError will wait for every given task to emit a result or
// return as soon as an error is encountered, stopping all other tasks.
func AllOrError(awaitables ...Awaitable) ([]interface{}, error) {
	defer stop.All(awaitableToStopable(awaitables...)...)

	errCh := make(chan error, 1)
	valueCh := make(chan interface{}, 1)

	for _, v := range awaitables {
		awaitable := v
		go func() {
			v, err := awaitable.Result()
			if err != nil {
				errCh <- errors.Wrap(err, 0)
				return
			}
			valueCh <- v
		}()
	}

	values := []interface{}{}
	for {
		select {
		case err := <-errCh:
			return nil, errors.Wrap(err, 0)
		case value := <-valueCh:
			values = append(values, value)
			if len(values) == len(awaitables) {
				return values, nil
			}
		}
	}
}

// AllOrErrorWithTimeout does the same as AllOrError but allows you to set a
// timeout for waiting for other tasks to stop.
func AllOrErrorWithTimeout(timeout time.Duration, awaitables ...Awaitable) ([]interface{}, error) {
	defer stop.AllWithTimeout(timeout, awaitableToStopableWithTimeout(awaitables...)...)

	errCh := make(chan error, 1)
	valueCh := make(chan interface{}, 1)

	for _, v := range awaitables {
		awaitable := v
		go func() {
			v, err := awaitable.Result()
			if err != nil {
				errCh <- errors.Wrap(err, 0)
				return
			}
			valueCh <- v
		}()
	}

	values := []interface{}{}
	for {
		select {
		case err := <-errCh:
			return nil, errors.Wrap(err, 0)
		case value := <-valueCh:
			values = append(values, value)
			if len(values) == len(awaitables) {
				return values, nil
			}
		}
	}
}

// FastAllOrError does the same as AllOrError but does not wait for all other
// tasks to stop, it does tell them to stop it just doesn't wait for them to stop.
func FastAllOrError(awaitables ...Awaitable) ([]interface{}, error) {
	defer stop.AllAsync(awaitableToStopable(awaitables...)...)

	errCh := make(chan error, 1)
	valueCh := make(chan interface{}, 1)

	for _, v := range awaitables {
		awaitable := v
		go func() {
			v, err := awaitable.Result()
			if err != nil {
				errCh <- errors.Wrap(err, 0)
				return
			}
			valueCh <- v
		}()
	}

	values := []interface{}{}
	for {
		select {
		case err := <-errCh:
			return nil, errors.Wrap(err, 0)
		case value := <-valueCh:
			values = append(values, value)
			if len(values) == len(awaitables) {
				return values, nil
			}
		}
	}
}

// Any will wait for the first task to emit a result (or an error)
// and return that, stopping all other tasks.
func Any(awaitables ...Awaitable) (interface{}, error) {
	defer stop.All(awaitableToStopable(awaitables...)...)

	valueCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)

	for _, v := range awaitables {
		awaitable := v
		go func() {
			v, e := awaitable.Result()
			if e != nil {
				errCh <- errors.Wrap(e, 0)
			}
			valueCh <- v
		}()
	}

	select {
	case v := <-valueCh:
		return v, nil
	case e := <-errCh:
		return nil, errors.Wrap(e, 0)
	}
}

// AnyWithTimeout does the same as Any but allows you to set a
// timeout for waiting for other tasks to stop.
func AnyWithTimeout(timeout time.Duration, awaitables ...Awaitable) (interface{}, error) {
	defer stop.AllWithTimeout(timeout, awaitableToStopableWithTimeout(awaitables...)...)

	valueCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)

	for _, v := range awaitables {
		awaitable := v
		go func() {
			v, e := awaitable.Result()
			if e != nil {
				errCh <- errors.Wrap(e, 0)
			}
			valueCh <- v
		}()
	}

	select {
	case v := <-valueCh:
		return v, nil
	case e := <-errCh:
		return nil, errors.Wrap(e, 0)
	}
}

// FastAny does the same as Any but does not wait for all other tasks to stop,
// it does tell them to stop it just doesn't wait for them to stop.
func FastAny(awaitables ...Awaitable) (interface{}, error) {
	defer stop.AllAsync(awaitableToStopable(awaitables...)...)

	valueCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)

	for _, v := range awaitables {
		awaitable := v
		go func() {
			v, e := awaitable.Result()
			if e != nil {
				errCh <- errors.Wrap(e, 0)
			}
			valueCh <- v
		}()
	}

	select {
	case v := <-valueCh:
		return v, nil
	case e := <-errCh:
		return nil, errors.Wrap(e, 0)
	}
}

// ErrTaskFailed is returned by the All methods when at least one task returns an error.
type ErrTaskFailed struct {
	Errors []error
}

func (e *ErrTaskFailed) Error() string {
	return "await: one or more errors were returned from the awaited tasks"
}

func awaitableToStopable(awaitables ...Awaitable) []stop.Stopable {
	stopables := []stop.Stopable{}
	for _, awaitable := range awaitables {
		v, ok := awaitable.(stop.Stopable)
		if ok {
			stopables = append(stopables, v)
		}
	}
	return stopables
}

func awaitableToStopableWithTimeout(awaitables ...Awaitable) []stop.StopableWithTimeout {
	stopables := []stop.StopableWithTimeout{}
	for _, awaitable := range awaitables {
		v, ok := awaitable.(stop.StopableWithTimeout)
		if ok {
			stopables = append(stopables, v)
		}
	}
	return stopables
}
