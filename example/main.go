package main

import (
	"fmt"
	"time"

	"github.com/brad-jones/goasync/await"
	"github.com/brad-jones/goasync/task"
	uuid "github.com/satori/go.uuid"
)

func randomNameAsync() *task.Task {
	return task.New(func(t *task.Internal) {
		time.Sleep(1 * time.Second)
		t.Resolve(uuid.NewV4().String())
	})
}

func sayHelloAsync(greeting *task.Task) *task.Task {
	return task.New(func(t *task.Internal) {
		v, e := greeting.Result()
		if e != nil {
			t.Reject(e)
		}
		time.Sleep(1 * time.Second)
		if v == "" {
			t.Reject(fmt.Errorf("greeting was not set"))
		} else {
			t.Resolve("hello: " + v.(string))
		}
	})
}

func main() {
	start := time.Now()
	fmt.Println("START", start)

	name1 := randomNameAsync()
	name2 := randomNameAsync()

	hello1 := sayHelloAsync(name1)
	hello2 := sayHelloAsync(name2)
	hello3 := sayHelloAsync(name2)

	fmt.Println("TASKS STARTED", time.Since(start))

	greetings, err := await.AllOrError(hello1, hello2, hello3)
	if err != nil {
		panic(err)
	}

	for _, greeting := range greetings {
		fmt.Println(greeting)
	}

	fmt.Println("END", time.Since(start))
}
