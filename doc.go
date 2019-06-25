/*
Package goasync is a helper framework for writing asynchronous code in go.
It's primary goal is to reduce the amount of boilier plate code one has to
write to do concurrent tasks.

Preface: My ideas presented here may not necessarily be theoretically perfect
but I am taking a pragmatic approach to my development of golang code, while at
the same time keeping in mind that go is go and not another language... go is
verbose suck it up and move on.

TLDR: Don't care about my journey through golang's concurrency model,
just want to know how this library works, skip down to The Task API.

Prior Art

https://github.com/chebyrash/promise

https://github.com/fanliao/go-promise

https://github.com/capitalone/go-future-context

https://github.com/rafaeldias/async

Other Reading

http://www.golangpatterns.info/concurrency/futures

https://stackoverflow.com/questions/35926173/implementing-promise-with-channels-in-go

https://medium.com/strava-engineering/futures-promises-in-the-land-of-golang-1453f4807945

https://www.reddit.com/r/golang/comments/3qhgsr/futurespromises_in_the_land_of_golang

Async Functions in Go

Sure you can prefix any function call with `go` and it will run asynchronously
but this isn't always the full story, more often than not you have to deal with
the results of that function call, including any error handling strategy and at
the very least some sort of method of ensuring it actually runs to completion
not to mention cancelation.

This is what I consider to be the basic async function:

	func fooAsync(bar <-chan string) (<-chan string, <-chan error) {
		resolver := make(chan string, 1)
		rejector := make(chan error, 1)
		go func() {
			v, err := doSomething(<-bar)
			if err != nil {
				rejector <- err
			} else {
				resolver <- v
			}
		}()
		return resolver, rejector
	}

Ok so lets pull that apart:

* Lets say it's idiomatic to add the `Async` suffix to any function name,
much like you might prefix a function name with `Must` to denote it might panic.
This follows other languages such as C#.

* Any inputs to an async function must be channels, this is so that when that
input is returned from another async function that channel can just be passed
in and read from the goroutine, not blocking anyone else.

* Similarly all returned values from an async function must be channels.
Effectively we have what might be similar to a Javascript promise but broken
into 2 parts, the resolver and the rejector.

* Resolvers and rejectors should be buffered such that they can be executed
before anything has been setup to read the channels.

* A resolver or rejector will only ever send a single value to the channel.

* A resolver or rejector can only be read once. A simple trick to "tee" a channel:

	ch2 := make(chan string, 1)
	ch3 := make(chan string, 1)
	go func() {
		v := <-ch1
		ch2 <- v
		ch3 <- v
	}()

Adding Cancelation

	func fooAsync(bar <-chan string) (resolver <-chan string, rejector <-chan error, stopper chan<- struct{}) {
		resolver = make(chan string, 1)
		rejector = make(chan error, 1)
		stopper = make(chan struct{}, 1)
		go func() {
			for {
				select {
				case <-stopper:
					return
				case <-time.After(5 * time.Second):
					resolver <- value
					return
				default:
					v, err := doSomething(<-bar)
					if err != nil {
						rejector <- err
					} else {
						value = value + v
					}
				}
			}
		}()
		return resolver, rejector, stopper
	}

* So we just written 24 lines of code and only one of them actually does anything
of any importance.

* While there are some cases where having the resolver, rejector & stopper
separate it gets hard to keep track with many variables. What if fooAsync
above wanted to see the error of bar, you would have to pass that in too.

* Why not use context.Context? Plenty of reading about that
https://dave.cheney.net/2017/08/20/context-isnt-for-cancellation

Awaiting in Go

Your friend is `select`, basically consider it be the replacement for the
keyword `await` used in other languages.

	resolver, rejector := barAsync()
	select {
	case v := <-resolver:
		fmt.Println("we got an value", v)
	case err := <-rejector:
		fmt.Println("we got an error", err)
	}

Await Any: This is how you might await for "any" of the results
from a collection of async calls.

	resolver1, rejector1 := fooAsync()
	resolver2, rejector2 := barAsync()
	resolver3, rejector3 := bazAsync()

	var value string
	var err error
	select {
	case v := <-resolver1:
		value = v
	case v := <-resolver2:
		value = v
	case v := <-resolver3:
		value = v
	case e := <-rejector1:
		err = e
	case e := <-rejector2:
		err = e
	case e := <-rejector3:
		err = e
	}

	if err != nil {
		panic(err)
	}
	fmt.Println("we got a value", value)

Await All: This is how you might await for "all" of the results
from a collection of async calls.

	resolver1, rejector1 := fooAsync()
	resolver2, rejector2 := barAsync()
	resolver3, rejector3 := bazAsync()

	values := []string{}
	errs := []error{}
	for i := 0; i < 3; i++ { // no this should not be 6
		select {
		case v := <-resolver1:
			values = append(values, v)
		case v := <-resolver2:
			values = append(values, v)
		case v := <-resolver3:
			values = append(values, v)
		case e := <-rejector1:
			errs = append(errs, e)
		case e := <-rejector2:
			errs = append(errs, e)
		case e := <-rejector3:
			errs = append(errs, e)
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println("oh no an error", err)
		}
		panic("we got errors")
	}
	for _, value := range values {
		fmt.Println("we got a value", value)
	}

Chaining Async Functions into a Pipeline

Normally you will want to chain additional actions to take place, in the
fastest way possible. Usually those functions will consume the outputs from
the previous functions. So you can just pass a channel into the next function
and it will wait on the channel inside it's own goroutine.

	aCh := make(chan string, 1)
	aCh <- "aFile"
	aTag, aErr := BuildDockerImageAsync(aCh)

	bCh := make(chan string, 1)
	bCh <- "bFile"
	bTag, bErr := BuildDockerImageAsync(bCh)

	cCh := make(chan string, 1)
	cCh <- "cFile"
	cTag, cErr := BuildDockerImageAsync(cCh)

	aPubDone, aPubErr := PublishDockerImageAsync(aTag)
	bPubDone, bPubErr := PublishDockerImageAsync(bTag)
	cPubDone, cPubErr := PublishDockerImageAsync(cTag)

	aDeployDone, aDeployErr := DeployDockerImageAsync(aPubDone)
	bDeployDone, bDeployErr := DeployDockerImageAsync(bPubDone)
	cDeployDone, cDeployErr := DeployDockerImageAsync(cPubDone)

	done := make(chan struct{}, 1)
	go func(){
		defer close(done)
		for i := 0; i < 3; i++ {
			select {
				case <-aDeployDone:
				case <-bDeployDone:
				case <-cDeployDone:
			}
		}
	}()

	var err error
	select {
	case <-done:
	case e := <-aErr:
		err = e
	case e := <-bErr:
		err = e
	case e := <-cErr:
		err = e
	case e := <-aPubErr:
		err = e
	case e := <-bPubErr:
		err = e
	case e := <-cPubErr:
		err = e
	case e := <-aDeployErr:
		err = e
	case e := <-bDeployErr:
		err = e
	case e := <-cDeployErr:
		err = e
	}
	if err != nil {
		panic(err)
	}

Just write a Synchronous API and let the consumer use it Asynchronously if they want

My issue with this is the boilerplate code one has to write to do this.
Whatever happened to DRY? I believe the methodologies outlined by this
library strike a reasonable balance between writing idiomatic go and going
insane repeating yourself everywhere.

Whats more this is just one way things can be done, if you need something more
powerful for a particular use case then the full power of go's concurrency model
is still there, this library doesn't take any of that away.

Up to this point all I have shown you is how to do some stuff with vanilla go
and I hope that not only does it illustrate some of my frustrations with the
language but also acts a useful reference to go back to when you need to do
something more complex.

The Task API

Here is an example of `fooAsync` but this time it uses https://github.com/brad-jones/goasync/task

	func fooAsync(bar *task.Task) *task.Task {
		return task.New(func(t *task.Internal){
			// normally you would call this at the start of any long running loop
			if t.ShouldStop() {
				return
			}

			res, err := bar.Result()
			if err != nil {
				t.Reject(e) // or you could do something based on the error
				return
			}

			v, err := intToString(res.(int))
			if err != nil {
				t.Reject(e)
			} else {
				t.Resolve(v)
			}
		})
	}

	t := fooAsync(task.Resolved(1))
	v, err := t.Result()
	castedV := v.(string)

The Await API

Tasks can be awaited using https://github.com/brad-jones/goasync/await

	values, errors := await.All(task1, task2, task3)
	values, error := await.AllOrError(task1, task2, task3)
	value, error := await.Any(task1, task2, task3)

The awaiters that return early (before all tasks are complete) such as Any will
cooperatively stop the remaining tasks. So cancelation will happen automatically.
If you do not care for cancelation and wish to have the awaiter return as soon
as possible you may uses the `Fast` awaiters.

`Fast` awaiters will still ask any remaining tasks to
cooperatively stop but they do so asynchronously.

	value, error := await.FastAny(task1, task2, task3)

Or perhaps you might like to use a timeout.

	value, error := await.AnyWithTimeout(5 * time.Second, task1, task2, task3)

Not Type Safe

Due to go's lack of generics the only sane why this package can be created is
by the use of the `interface{}` type. This means that all values that are
returned from a task's `Result()` method or an awaiter must be casted correctly
by the caller.
*/
package goasync
