/*
Package goasync is a helper framework for writing asynchronous code in go.
It's 2 primary goals are to remain type safe while reducing boilier plate code.

Preface: My ideas presented here may not necessarily be theoretically perfect
but I am taking a pragmatic approach to my development of golang code, while at
the same time keeping in mind that go is go and not another language... go is
verbose suck it up and move on.

Prior Art

The basic issue with most of these solutions is they end up using
`interface{}` and/or reflection. Both of which I would love to avoid.

https://github.com/chebyrash/promise

https://github.com/fanliao/go-promise


Other Reading

http://www.golangpatterns.info/concurrency/futures

https://stackoverflow.com/questions/35926173/implementing-promise-with-channels-in-go

https://medium.com/strava-engineering/futures-promises-in-the-land-of-golang-1453f4807945

https://www.reddit.com/r/golang/comments/3qhgsr/futurespromises_in_the_land_of_golang

Async Functions in Go

Sure you can prefix any function call with `go` and it will run asynchronously
but this isn't always the full story, more often than not you have to deal with
the results of that function call, including any error handling strategy and at
the very least some sort of method of ensuring it actually runs to completion.

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

Ok so lets pull that appart:

* Lets say it's idiomatic to add the `Async` suffix to any function name,
much like you might prefix a function name with `Must` to denote it might panic.
This follows other langauges such as C#.

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

Just write a Synchronous API and let the consumer use it Asynchronously if they want

My issue with this is the boilerplate code one has to write to do this.
Whatever happened to DRY? I believe the methodologies outlined by this
library strike a reasonable balance between writing idiomatic go and going
insane repeating yourself everywhere.

Whats more this is just one way things can be done, if you need something more
powerful for a particular use case then the full power of go's concurrency model
is still there, this library doesn't take any of that away.

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

Making Channels

Up to this point all I have shown you is how to do some stuff with vanilla go.
Now I will introduce the first set of helper functions that this library provides.

Due to the rule above that, "any inputs to an async function must be channels",
it becomes a pain to start off an asynchronous chain like the above.

Now if you know for certian that a particular function will always be used at
the start and never, ever have to accept input from another async function then
by all means just make the function accept the value directly and move on.

However for the cases that you can not be so certian how your async function
will be used then you can use the `MakeCh` helpers like so:

	import "github.com/brad-jones/goasync/ch"

	aTag, aErr := BuildDockerImageAsync(ch.MakeString("aFile"))
	bTag, bErr := BuildDockerImageAsync(ch.MakeString("bFile"))
	cTag, cErr := BuildDockerImageAsync(ch.MakeString("cFile"))

Await Helpers

Using select to await for things is very powerful but it is also very verbose
so this library provides a collection of await helpers, for example:

	import "github.com/brad-jones/goasync/await"

	firstString := await.AnyString(aTag, bTag, cTag)
	allStrings := await.AllStrings(aTag, bTag, cTag)

Async Awaiters

Yes you heard me, we can have async awaiters, an example is worth a 1000 words:

	import "github.com/brad-jones/goasync/await"

	aDeployDone, aDeployErr := DeployDockerImageAsync(aPubDone)
	bDeployDone, bDeployErr := DeployDockerImageAsync(bPubDone)
	cDeployDone, cDeployErr := DeployDockerImageAsync(cPubDone)

	select {
	case <-await.AllAsync(aDeployDone, bDeployDone, cDeployDone):
	case err := <-await.AnyErrorAsync(aDeployErr, bDeployErr, cDeployErr):
		panic(err)
	}

Cancellations

If you wish for your async function to return early in the event another async
function returns an error, or the context is otherwise canceled, you might do
something like:

	// a so called long running function will usually involve a loop
	// of some sort that allows us to check the context on each iteration
	func longRunningAsyncWithContext(ctx context.Context, bar <-chan []string) (<-chan []string, <-chan error) {
		resolver := make(chan []string, 1)
		rejector := make(chan error, 1)
		go func() {
			output := []string{}
			for _, b := range <-bar {
				select {
				case <-ctx.Done():
					return
				default:
					v, err := doSomething(b)
					if err != nil {
						rejector <- err
						return
					}
					output = append(output, v)
				}
			}
			resolver <- output
		}()
		return resolver, rejector
	}

	// a so called short running function can only check the context once at the start
	// so if this function has already started and the cancel function is invoked
	// it will have to run to completion, unless of course main() exits.
	func shortRunningAsyncWithContext(ctx context.Context, bar <-chan string) (<-chan string, <-chan error) {
		resolver := make(chan []string, 1)
		rejector := make(chan error, 1)
		go func() {
			select {
			case <-ctx.Done():
				return
			default:
				v, err := doSomething(b)
				if err != nil {
					rejector <- err
				} else {
					resolver <- v
				}
			}
		}()
		return resolver, rejector
	}

	func main() {
		ctx, cancel := context.WithCancel(context.Background())
		resolver1, rejector1 := longRunningAsyncWithContext(ctx, ch.MakeStringSlice([]string{"a","b","c"}))
		resolver2, rejector2 := shortRunningAsyncWithContext(ctx, ch.MakeString("a"))

		select {
		case <-await.AllAsync(resolver1, resolver2):
		case err := <-await.AnyErrorAsync(rejector1, rejector2):
			cancel()
		}
	}

Using Genny to generate helpers for your own types

Genny (https://github.com/cheekybits/genny) is a code-generation generics solution.
This library uses it to generate all the type safe variations of our helpers.
If you wanted to generate your own custom helpers you might do something like:

	go get -u github.com/cheekybits/genny && \
	curl -sL "https://github.com/brad-jones/goasync/raw/master/await/generic.genny" | \
	genny -pkg="mypkg" gen "Awaitable=*MyType" > ./awaiters-mytype.go
*/
package goasync
