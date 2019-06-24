package await

//go:generate genny -in=generic.genny -out=typed.go gen "Awaitable=BUILTINS"

import (
	"sync"

	"github.com/brad-jones/goasync/stop"
	"github.com/brad-jones/goasync/task"
	"github.com/brad-jones/goerr"
)

type Awaitable interface {
	Result() (interface{}, error)
}
type AwaitableBool interface {
	Result() (bool, error)
}
type AwaitableByte interface {
	Result() (byte, error)
}
type AwaitableComplex128 interface {
	Result() (complex128, error)
}
type AwaitableComplex64 interface {
	Result() (complex64, error)
}
type AwaitableError interface {
	Result() (error, error)
}
type AwaitableFloat32 interface {
	Result() (float32, error)
}
type AwaitableFloat64 interface {
	Result() (float64, error)
}
type AwaitableInt interface {
	Result() (int, error)
}
type AwaitableInt16 interface {
	Result() (int16, error)
}
type AwaitableInt32 interface {
	Result() (int32, error)
}
type AwaitableInt64 interface {
	Result() (int64, error)
}
type AwaitableInt8 interface {
	Result() (int8, error)
}
type AwaitableRune interface {
	Result() (rune, error)
}
type AwaitableString interface {
	Result() (string, error)
}
type AwaitableUint interface {
	Result() (uint, error)
}
type AwaitableUint16 interface {
	Result() (uint16, error)
}
type AwaitableUint32 interface {
	Result() (uint32, error)
}
type AwaitableUint64 interface {
	Result() (uint64, error)
}
type AwaitableUint8 interface {
	Result() (uint8, error)
}
type AwaitableUintptr interface {
	Result() (uintptr, error)
}
type AwaitableBoolSlice interface {
	Result() ([]bool, error)
}
type AwaitableByteSlice interface {
	Result() ([]byte, error)
}
type AwaitableComplex128Slice interface {
	Result() ([]complex128, error)
}
type AwaitableComplex64Slice interface {
	Result() ([]complex64, error)
}
type AwaitableErrorSlice interface {
	Result() ([]error, error)
}
type AwaitableFloat32Slice interface {
	Result() ([]float32, error)
}
type AwaitableFloat64Slice interface {
	Result() ([]float64, error)
}
type AwaitableIntSlice interface {
	Result() ([]int, error)
}
type AwaitableInt16Slice interface {
	Result() ([]int16, error)
}
type AwaitableInt32Slice interface {
	Result() ([]int32, error)
}
type AwaitableInt64Slice interface {
	Result() ([]int64, error)
}
type AwaitableInt8Slice interface {
	Result() ([]int8, error)
}
type AwaitableRuneSlice interface {
	Result() ([]rune, error)
}
type AwaitableStringSlice interface {
	Result() ([]string, error)
}
type AwaitableUintSlice interface {
	Result() ([]uint, error)
}
type AwaitableUint16Slice interface {
	Result() ([]uint16, error)
}
type AwaitableUint32Slice interface {
	Result() ([]uint32, error)
}
type AwaitableUint64Slice interface {
	Result() ([]uint64, error)
}
type AwaitableUint8Slice interface {
	Result() ([]uint8, error)
}
type AwaitableUintptrSlice interface {
	Result() ([]uintptr, error)
}

// All will wait for every given task to emit a result.
// The results (& errors) will be returned in a slice
// ordered the same as the input.
//
// Not type safe, use with care.
func All(awaitables ...interface{}) ([]interface{}, error) {
	awaited := []interface{}{}
	errors := []error{}

	for _, awaitable := range awaitables {
		var v interface{}
		var e error

		switch awaitable := awaitable.(type) {
		case Awaitable:
			v, e = awaitable.Result()
		case AwaitableBool:
			v, e = awaitable.Result()
		case AwaitableByte:
			v, e = awaitable.Result()
		case AwaitableComplex128:
			v, e = awaitable.Result()
		case AwaitableComplex64:
			v, e = awaitable.Result()
		case AwaitableError:
			v, e = awaitable.Result()
		case AwaitableFloat32:
			v, e = awaitable.Result()
		case AwaitableFloat64:
			v, e = awaitable.Result()
		case AwaitableInt:
			v, e = awaitable.Result()
		case AwaitableInt16:
			v, e = awaitable.Result()
		case AwaitableInt32:
			v, e = awaitable.Result()
		case AwaitableInt64:
			v, e = awaitable.Result()
		case AwaitableInt8:
			v, e = awaitable.Result()
		case AwaitableRune:
			v, e = awaitable.Result()
		case AwaitableString:
			v, e = awaitable.Result()
		case AwaitableUint:
			v, e = awaitable.Result()
		case AwaitableUint16:
			v, e = awaitable.Result()
		case AwaitableUint32:
			v, e = awaitable.Result()
		case AwaitableUint64:
			v, e = awaitable.Result()
		case AwaitableUint8:
			v, e = awaitable.Result()
		case AwaitableUintptr:
			v, e = awaitable.Result()
		case AwaitableBoolSlice:
			v, e = awaitable.Result()
		case AwaitableByteSlice:
			v, e = awaitable.Result()
		case AwaitableComplex128Slice:
			v, e = awaitable.Result()
		case AwaitableComplex64Slice:
			v, e = awaitable.Result()
		case AwaitableErrorSlice:
			v, e = awaitable.Result()
		case AwaitableFloat32Slice:
			v, e = awaitable.Result()
		case AwaitableFloat64Slice:
			v, e = awaitable.Result()
		case AwaitableIntSlice:
			v, e = awaitable.Result()
		case AwaitableInt16Slice:
			v, e = awaitable.Result()
		case AwaitableInt32Slice:
			v, e = awaitable.Result()
		case AwaitableInt64Slice:
			v, e = awaitable.Result()
		case AwaitableInt8Slice:
			v, e = awaitable.Result()
		case AwaitableRuneSlice:
			v, e = awaitable.Result()
		case AwaitableStringSlice:
			v, e = awaitable.Result()
		case AwaitableUintSlice:
			v, e = awaitable.Result()
		case AwaitableUint16Slice:
			v, e = awaitable.Result()
		case AwaitableUint32Slice:
			v, e = awaitable.Result()
		case AwaitableUint64Slice:
			v, e = awaitable.Result()
		case AwaitableUint8Slice:
			v, e = awaitable.Result()
		case AwaitableUintptrSlice:
			v, e = awaitable.Result()
		default:
			panic("TODO: use reflection here")
		}

		awaited = append(awaited, v)
		errors = append(errors, e)
	}

	if len(errors) > 0 {
		return nil, &ErrTaskFailed{
			Errors: errors,
		}
	}
	return awaited, nil
}

// AllAsync is an asynchronous version of All.
//
// Not type safe, use with care.
func AllAsync(awaitables ...interface{}) *task.Task {
	return task.New(func(t *task.TaskInternal) {
		v, e := All(awaitables...)
		goerr.Check(e)
		t.Resolve(v)
	})
}

// Any will wait for the first task to emit a result (or an error),
// return that and stop all other tasks.
//
// Not type safe, use with care.
func Any(awaitables ...interface{}) (interface{}, error) {
	defer stop.All(stop.ToStopables(awaitables...)...)

	var value interface{}
	var err error
	type doneOnce struct {
		o  sync.Once
		ch chan struct{}
	}
	done := doneOnce{
		ch: make(chan struct{}, 1),
	}
	closeDone := func() {
		done.o.Do(func() {
			close(done.ch)
		})
	}

	for _, awaitable := range awaitables {
		switch awaitable := awaitable.(type) {
		case Awaitable:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableBool:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableByte:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableComplex128:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableComplex64:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableError:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableFloat32:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableFloat64:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableInt:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableInt16:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableInt32:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableInt64:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableInt8:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableRune:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableString:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableUint:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableUint16:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableUint32:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableUint64:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableUint8:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableUintptr:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableBoolSlice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableByteSlice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableComplex128Slice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableComplex64Slice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableErrorSlice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableFloat32Slice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableFloat64Slice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableIntSlice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableInt16Slice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableInt32Slice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableInt64Slice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableInt8Slice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableRuneSlice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableStringSlice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableUintSlice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableUint16Slice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableUint32Slice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableUint64Slice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableUint8Slice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		case AwaitableUintptrSlice:
			go func() {
				defer closeDone()
				value, err = awaitable.Result()
			}()
		default:
			panic("TODO: use reflection here")
		}
	}
	<-done.ch
	return value, err
}

// AnyAsync is an asynchronous version of Any.
//
// Not type safe, use with care.
func AnyAsync(awaitables ...interface{}) *task.Task {
	return task.New(func(t *task.TaskInternal) {
		v, e := Any(awaitables...)
		goerr.Check(e)
		t.Resolve(v)
	})
}
