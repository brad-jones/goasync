package main

import (
	"fmt"
	"time"

	"github.com/brad-jones/goasync/v2/task"
)

func cancelableAsync() *task.Task {
	return task.New(func(t *task.Internal) {
		chainedTask := chainedCancelableAsync()
		chainedTask.Stopper = t.Stopper
		for i := 1; i < 10; i++ {
			if t.ShouldStop() {
				fmt.Println("cancelableAsync: I stopped cooperatively")
				t.Reject(fmt.Errorf("cancelableAsync took too long"))
				return
			}
			fmt.Println("cancelableAsync: running for the", i, "time")
			time.Sleep(200 * time.Millisecond)
		}
		chainedTask.Result()
		fmt.Println("cancelableAsync: finished work")
	})
}

func chainedCancelableAsync() *task.Task {
	return task.New(func(t *task.Internal) {
		for i := 1; i < 10; i++ {
			if t.ShouldStop() {
				fmt.Println("chainedCancelableAsync: I stopped cooperatively")
				t.Reject(fmt.Errorf("chainedCancelableAsync took too long"))
				return
			}
			fmt.Println("chainedCancelableAsync: running for the", i, "time")
			time.Sleep(200 * time.Millisecond)
		}
		fmt.Println("chainedCancelableAsync: finished work")
	})
}

func main() {
	start := time.Now()
	fmt.Println("START", start)

	if _, err := cancelableAsync().ResultWithTimeout(1*time.Second, 210*time.Millisecond); err != nil {
		fmt.Println(err)
	}

	fmt.Println("END", time.Since(start))
}
