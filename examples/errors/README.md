# Error Handling

goasync integrates with [goerr](https://github.com/brad-jones/goerr), which
essentially stores a single stack frame against a wrapped error value.
By unwrapping the entire error chain one can then generate a pseudo callstack.

This has the advantage of working across goroutine boundaries, even for
panics (that are recovered and turned into an error value by goasync).

The intent is not that we capture the exact state of the stack when an error
happens, including every function call. For a library that does that,
see <https://github.com/go-errors/errors>.

The intent here is to attach relevant contextual information (messages,
variables) at strategic places along the call stack, keeping stack traces
compact and maximally useful.

## Expected Output

```
something went wrong

main.crash3.func1:C:/Users/brad.jones/Projects/Personal/goasync/examples/errors/main.go:32
        t.Reject(goerr.New("something went wrong"))
main.crash2.func1:C:/Users/brad.jones/Projects/Personal/goasync/examples/errors/main.go:23
        t.Reject(e)
main.crash1.func1:C:/Users/brad.jones/Projects/Personal/goasync/examples/errors/main.go:12
        t.Reject(e)
```
