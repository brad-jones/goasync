package main

import (
	"fmt"
	"time"

	"github.com/brad-jones/goasync/v2/await"
	"github.com/brad-jones/goasync/v2/task"
)

func foo() *task.Task {
	return task.New(func(t *task.Internal) {
		time.Sleep(1 * time.Second)
		t.Resolve("foo did some work")
	})
}

func bar() *task.Task {
	return task.New(func(t *task.Internal) {
		time.Sleep(2 * time.Second)
		t.Resolve("bar did some work")
	})
}

func main() {
	start := time.Now()
	fmt.Println("START", start)

	s := await.Stream(foo(), bar())
	for s.Wait() {
		fmt.Println(s.MustResult().(string))
	}

	fmt.Println("END", time.Since(start))
}
