package goasync

import (
	"context"
)

// ShouldCancel is non blocking and will return true if the context's
// Done channel has been closed.
//
// 	func fooAsync(bar <-chan string, ctx ...context.Context) (<-chan string, <-chan error) {
// 		resolver := make(chan string, 1)
// 		rejector := make(chan error, 1)
// 		go func() {
// 			if goasync.ShouldCancel(ctx) {
// 				return
// 			}
// 			v, err := doSomething(<-bar)
// 			if err != nil {
// 				rejector <- err
// 			} else {
// 				resolver <- v
// 			}
// 		}()
// 		return resolver, rejector
// 	}
//
// I know the context should always come first but again I am being pragmatic
// and saving myself from defining another method that accepts a context.
//
// I guess this won't always work, what if the function actually used variadic
// arguments, then yeah you would have to create a "WithContext" function.
//
// I guess you could still have one function and pass in nil for context but
// that kind of sucks too.
func ShouldCancel(ctx []context.Context) bool {
	if len(ctx) == 1 {
		select {
		default:
		case <-ctx[0].Done():
			return true
		}
	}
	return false
}
