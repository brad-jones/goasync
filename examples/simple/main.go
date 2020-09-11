package main

import (
	"fmt"
	"time"

	"github.com/brad-jones/goasync/v2/await"
	"github.com/brad-jones/goasync/v2/task"
)

func doSomeWorkAsync() *task.Task {
	return task.New(func(t *task.Internal) {
		time.Sleep(1 * time.Second)
		fmt.Println("doing work")
	})
}

func main() {
	start := time.Now()
	fmt.Println("START", start)

	task1 := doSomeWorkAsync()
	task2 := doSomeWorkAsync()
	await.All(task1, task2)

	fmt.Println("END", time.Since(start))
}
