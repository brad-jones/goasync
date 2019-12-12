# goasync

[![GoReport](https://goreportcard.com/badge/brad-jones/goasync)](https://goreportcard.com/report/brad-jones/goasync)
[![GoLang](https://img.shields.io/badge/golang-%3E%3D%201.13.4-lightblue.svg)](https://golang.org)
[![GoDoc](https://godoc.org/github.com/brad-jones/goasync?status.svg)](https://godoc.org/github.com/brad-jones/goasync)
[![License](https://img.shields.io/github/license/brad-jones/goasync.svg)](https://github.com/brad-jones/goasync/blob/master/LICENSE)

Package goasync is a helper framework for writing asynchronous code in go.
It's primary goal is to reduce the amount of boilier plate code one has to
write to do concurrent tasks.

## Usage

`go get -u github.com/brad-jones/goasync/...`

```go
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

```
