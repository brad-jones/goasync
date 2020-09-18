// Package await contains await helper functions for use with tasks.
package await

import (
	"time"

	"github.com/brad-jones/goasync/v2/stop"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
)

// All will wait for every given task to emit a result, the results (& errors)
// will be returned in a slice ordered the same as the input.
func All(awaitables ...*task.Task) ([]interface{}, error) {
	awaited := []interface{}{}
	awaitedErrors := []error{}

	for _, awaitable := range awaitables {
		v, e := awaitable.Result()
		if e != nil {
			awaitedErrors = append(awaitedErrors, goerr.Wrap(e))
		}
		awaited = append(awaited, v)
	}

	if len(awaitedErrors) > 0 {
		return nil, goerr.Wrap(&ErrTaskFailed{
			Errors: awaitedErrors,
		})
	}

	return awaited, nil
}

// MustAll does the same thing as All but panics if an error is encountered
func MustAll(awaitables ...*task.Task) []interface{} {
	v, e := All(awaitables...)
	goerr.Check(e)
	return v
}

// AllAsync does the same thing as All but does so asynchronously
func AllAsync(awaitables ...*task.Task) *task.Task {
	return task.New(func(t *task.Internal) {
		t.Resolve(MustAll(awaitables...))
	})
}

// AllOrError will wait for every given task to emit a result or
// return as soon as an error is encountered, stopping all other tasks.
func AllOrError(awaitables ...*task.Task) ([]interface{}, error) {
	defer stop.All(awaitables...)

	errCh := make(chan error, 1)
	valueCh := make(chan map[*task.Task]interface{}, 1)

	for _, v := range awaitables {
		awaitable := v
		go func() {
			v, err := awaitable.Result()
			if err != nil {
				errCh <- goerr.Wrap(err)
				return
			}
			valueCh <- map[*task.Task]interface{}{
				awaitable: v,
			}
		}()
	}

	values := map[*task.Task]interface{}{}
	for {
		select {
		case err := <-errCh:
			return nil, goerr.Wrap(err)
		case value := <-valueCh:
			for k, v := range value {
				values[k] = v
			}
			if len(values) == len(awaitables) {
				sortedValues := []interface{}{}
				for _, awaitable := range awaitables {
					sortedValues = append(sortedValues, values[awaitable])
				}
				return sortedValues, nil
			}
		}
	}
}

// MustAllOrError does the same thing as AllOrError but panics if an error is encountered
func MustAllOrError(awaitables ...*task.Task) []interface{} {
	v, e := AllOrError(awaitables...)
	goerr.Check(e)
	return v
}

// AllOrErrorAsync does the same thing as AllOrError but does so asynchronously
func AllOrErrorAsync(awaitables ...*task.Task) *task.Task {
	return task.New(func(t *task.Internal) {
		t.Resolve(MustAllOrError(awaitables...))
	})
}

// AllOrErrorWithTimeout does the same as AllOrError but allows you to set a
// timeout for waiting for other tasks to stop.
func AllOrErrorWithTimeout(timeout time.Duration, awaitables ...*task.Task) ([]interface{}, error) {
	defer stop.AllWithTimeout(timeout, awaitables...)

	errCh := make(chan error, 1)
	valueCh := make(chan map[*task.Task]interface{}, 1)

	for _, v := range awaitables {
		awaitable := v
		go func() {
			v, err := awaitable.Result()
			if err != nil {
				errCh <- goerr.Wrap(err)
				return
			}
			valueCh <- map[*task.Task]interface{}{
				awaitable: v,
			}
		}()
	}

	values := map[*task.Task]interface{}{}
	for {
		select {
		case err := <-errCh:
			return nil, goerr.Wrap(err)
		case value := <-valueCh:
			for k, v := range value {
				values[k] = v
			}
			if len(values) == len(awaitables) {
				sortedValues := []interface{}{}
				for _, awaitable := range awaitables {
					sortedValues = append(sortedValues, values[awaitable])
				}
				return sortedValues, nil
			}
		}
	}
}

// MustAllOrErrorWithTimeout does the same thing as AllOrErrorWithTimeout but panics if an error is encountered
func MustAllOrErrorWithTimeout(timeout time.Duration, awaitables ...*task.Task) []interface{} {
	v, e := AllOrErrorWithTimeout(timeout, awaitables...)
	goerr.Check(e)
	return v
}

// AllOrErrorWithTimeoutAsync does the same thing as AllOrErrorWithTimeout but does so asynchronously
func AllOrErrorWithTimeoutAsync(timeout time.Duration, awaitables ...*task.Task) *task.Task {
	return task.New(func(t *task.Internal) {
		t.Resolve(MustAllOrErrorWithTimeout(timeout, awaitables...))
	})
}

// FastAllOrError does the same as AllOrError but does not wait for all other
// tasks to stop, it does tell them to stop it just doesn't wait for them to stop.
func FastAllOrError(awaitables ...*task.Task) ([]interface{}, error) {
	defer stop.AllAsync(awaitables...)

	errCh := make(chan error, 1)
	valueCh := make(chan map[*task.Task]interface{}, 1)

	for _, v := range awaitables {
		awaitable := v
		go func() {
			v, err := awaitable.Result()
			if err != nil {
				errCh <- goerr.Wrap(err)
				return
			}
			valueCh <- map[*task.Task]interface{}{
				awaitable: v,
			}
		}()
	}

	values := map[*task.Task]interface{}{}
	for {
		select {
		case err := <-errCh:
			return nil, goerr.Wrap(err)
		case value := <-valueCh:
			for k, v := range value {
				values[k] = v
			}
			if len(values) == len(awaitables) {
				sortedValues := []interface{}{}
				for _, awaitable := range awaitables {
					sortedValues = append(sortedValues, values[awaitable])
				}
				return sortedValues, nil
			}
		}
	}
}

// MustFastAllOrError does the same thing as FastAllOrError but panics if an error is encountered
func MustFastAllOrError(awaitables ...*task.Task) []interface{} {
	v, e := FastAllOrError(awaitables...)
	goerr.Check(e)
	return v
}

// FastAllOrErrorAsync does the same thing as FastAllOrError but does so asynchronously
func FastAllOrErrorAsync(awaitables ...*task.Task) *task.Task {
	return task.New(func(t *task.Internal) {
		t.Resolve(MustFastAllOrError(awaitables...))
	})
}

// Any will wait for the first task to emit a result (or an error)
// and return that, stopping all other tasks.
func Any(awaitables ...*task.Task) (interface{}, error) {
	defer stop.All(awaitables...)

	valueCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)

	for _, v := range awaitables {
		awaitable := v
		go func() {
			v, e := awaitable.Result()
			if e != nil {
				errCh <- goerr.Wrap(e)
			}
			valueCh <- v
		}()
	}

	select {
	case v := <-valueCh:
		return v, nil
	case e := <-errCh:
		return nil, goerr.Wrap(e)
	}
}

// MustAny does the same thing as Any but panics if an error is encountered
func MustAny(awaitables ...*task.Task) interface{} {
	v, e := Any(awaitables...)
	goerr.Check(e)
	return v
}

// AnyAsync does the same thing as Any but does so asynchronously
func AnyAsync(awaitables ...*task.Task) *task.Task {
	return task.New(func(t *task.Internal) {
		t.Resolve(MustAny(awaitables...))
	})
}

// AnyWithTimeout does the same as Any but allows you to set a
// timeout for waiting for other tasks to stop.
func AnyWithTimeout(timeout time.Duration, awaitables ...*task.Task) (interface{}, error) {
	defer stop.AllWithTimeout(timeout, awaitables...)

	valueCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)

	for _, v := range awaitables {
		awaitable := v
		go func() {
			v, e := awaitable.Result()
			if e != nil {
				errCh <- goerr.Wrap(e)
			}
			valueCh <- v
		}()
	}

	select {
	case v := <-valueCh:
		return v, nil
	case e := <-errCh:
		return nil, goerr.Wrap(e)
	}
}

// MustAnyWithTimeout does the same thing as AnyWithTimeout but panics if an error is encountered
func MustAnyWithTimeout(timeout time.Duration, awaitables ...*task.Task) interface{} {
	v, e := AnyWithTimeout(timeout, awaitables...)
	goerr.Check(e)
	return v
}

// AnyWithTimeoutAsync does the same thing as AnyWithTimeout but does so asynchronously
func AnyWithTimeoutAsync(timeout time.Duration, awaitables ...*task.Task) *task.Task {
	return task.New(func(t *task.Internal) {
		t.Resolve(MustAnyWithTimeout(timeout, awaitables...))
	})
}

// FastAny does the same as Any but does not wait for all other tasks to stop,
// it does tell them to stop it just doesn't wait for them to stop.
func FastAny(awaitables ...*task.Task) (interface{}, error) {
	defer stop.AllAsync(awaitables...)

	valueCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)

	for _, v := range awaitables {
		awaitable := v
		go func() {
			v, e := awaitable.Result()
			if e != nil {
				errCh <- goerr.Wrap(e)
			}
			valueCh <- v
		}()
	}

	select {
	case v := <-valueCh:
		return v, nil
	case e := <-errCh:
		return nil, goerr.Wrap(e)
	}
}

// MustFastAny does the same thing as FastAny but panics if an error is encountered
func MustFastAny(awaitables ...*task.Task) interface{} {
	v, e := FastAny(awaitables...)
	goerr.Check(e)
	return v
}

// FastAnyAsync does the same thing as FastAny but does so asynchronously
func FastAnyAsync(awaitables ...*task.Task) *task.Task {
	return task.New(func(t *task.Internal) {
		t.Resolve(MustFastAny(awaitables...))
	})
}

// ErrTaskFailed is returned by the All methods when at least one task returns an error.
type ErrTaskFailed struct {
	Errors []error
}

func (e *ErrTaskFailed) Error() string {
	return "await: one or more errors were returned from the awaited tasks"
}
