package main

import (
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
)

func crash1() *task.Task {
	return task.New(func(t *task.Internal) {
		v, e := crash2().Result()
		if e != nil {
			t.Reject(e)
			return
		}
		t.Resolve(v)
	})
}

func crash2() *task.Task {
	return task.New(func(t *task.Internal) {
		v, e := crash3().Result()
		if e != nil {
			t.Reject(e)
			return
		}
		t.Resolve(v)
	})
}

func crash3() *task.Task {
	return task.New(func(t *task.Internal) {
		t.Reject(goerr.New("something went wrong"))
	})
}

func main() {
	if _, err := crash1().Result(); err != nil {
		goerr.PrintTrace(err)
	}
}
