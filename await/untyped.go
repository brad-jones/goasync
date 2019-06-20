package await

//go:generate genny -in=generic.genny -out=typed.go gen "Awaitable=BUILTINS"

// All will wait for every given channel to emit a single result.
// The results will be returned in a slice ordered the same as the input channels.
//
// Not type safe, use with care.
func All(awaitables ...interface{}) []interface{} {
	awaited := []interface{}{}
	for _, awaitable := range awaitables {
		switch v := awaitable.(type) {
		case <-chan struct{}:
			awaited = append(awaited, <-v)
		case <-chan bool:
			awaited = append(awaited, <-v)
		case <-chan byte:
			awaited = append(awaited, <-v)
		case <-chan complex128:
			awaited = append(awaited, <-v)
		case <-chan complex64:
			awaited = append(awaited, <-v)
		case <-chan error:
			awaited = append(awaited, <-v)
		case <-chan float32:
			awaited = append(awaited, <-v)
		case <-chan float64:
			awaited = append(awaited, <-v)
		case <-chan int:
			awaited = append(awaited, <-v)
		case <-chan int16:
			awaited = append(awaited, <-v)
		case <-chan int64:
			awaited = append(awaited, <-v)
		case <-chan int8:
			awaited = append(awaited, <-v)
		case <-chan rune:
			awaited = append(awaited, <-v)
		case <-chan string:
			awaited = append(awaited, <-v)
		case <-chan uint:
			awaited = append(awaited, <-v)
		case <-chan uint16:
			awaited = append(awaited, <-v)
		case <-chan uint32:
			awaited = append(awaited, <-v)
		case <-chan uint64:
			awaited = append(awaited, <-v)
		case <-chan uintptr:
			awaited = append(awaited, <-v)
		default:
			awaited = append(awaited, <-awaitable.(<-chan interface{}))
		}
	}
	return awaited
}

// AllAsync is an asynchronous version of AwaitAll.
//
// Not type safe, use with care.
func AllAsync(awaitables ...interface{}) <-chan []interface{} {
	resolver := make(chan []interface{}, 1)
	go func() {
		resolver <- All(awaitables...)
	}()
	return resolver
}

// Any will wait for the first channel to emit a single result
// and return that, ignoring all other channels.
//
// Not type safe, use with care.
func Any(awaitables ...interface{}) interface{} {
	for {
		for _, awaitable := range awaitables {
			switch v := awaitable.(type) {
			case <-chan struct{}:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan bool:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan byte:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan complex128:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan complex64:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan error:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan float32:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan float64:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan int:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan int16:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan int64:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan int8:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan rune:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan string:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan uint:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan uint16:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan uint32:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan uint64:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			case <-chan uintptr:
				select {
				default:
				case awaited := <-v:
					return awaited
				}
			default:
				select {
				default:
				case awaited := <-awaitable.(<-chan interface{}):
					return awaited
				}
			}
		}
	}
}

// AnyAsync is an asynchronous version of AwaitAny.
//
// Not type safe, use with care.
func AnyAsync(awaitables ...interface{}) <-chan interface{} {
	resolver := make(chan interface{}, 1)
	go func() {
		resolver <- Any(awaitables...)
	}()
	return resolver
}
