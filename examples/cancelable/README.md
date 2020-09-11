# Cancelable Tasks

This example shows how you can cancel long running tasks, including chained
tasks by propagating the `Stopper` channel.

## Expected Output

```
START 2020-09-11 18:53:19.6311565 +1000 AEST m=+0.003013801
cancelableAsync: running for the 1 time
chainedCancelableAsync: running for the 1 time
chainedCancelableAsync: running for the 2 time
cancelableAsync: running for the 2 time
cancelableAsync: running for the 3 time
chainedCancelableAsync: running for the 3 time
chainedCancelableAsync: running for the 4 time
cancelableAsync: running for the 4 time
cancelableAsync: running for the 5 time
chainedCancelableAsync: running for the 5 time
chainedCancelableAsync: running for the 6 time
cancelableAsync: running for the 6 time
cancelableAsync: I stopped cooperatively
chainedCancelableAsync: I stopped cooperatively
result took too long to be returned: cancelableAsync took too long
END 1.2059064s
```
