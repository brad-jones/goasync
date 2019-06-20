# goasync

[![GoReport](https://goreportcard.com/badge/brad-jones/goasync)](https://goreportcard.com/report/brad-jones/goasync)
[![GoLang](https://img.shields.io/badge/golang-%3E%3D%201.12.6-lightblue.svg)](https://golang.org)
[![GoDoc](https://godoc.org/github.com/brad-jones/goasync?status.svg)](https://godoc.org/github.com/brad-jones/goasync)
[![License](https://img.shields.io/github/license/brad-jones/goasync.svg)](https://github.com/brad-jones/goasync/blob/master/LICENSE)

A helper framework for writing asynchronous code in go.
It's 2 primary goals are to remain type safe while reducing boilier plate code.

## Usage

`go get -u github.com/brad-jones/goasync/...`

```go
package main

import (
	"fmt"
	"time"

	"github.com/brad-jones/goasync/await"
	uuid "github.com/satori/go.uuid"
)

func randomNameAsync() (<-chan string, <-chan error) {
	resolver := make(chan string, 1)
	rejector := make(chan error, 1)
	go func() {
		time.Sleep(1 * time.Second)
		resolver <- uuid.NewV4().String()
	}()
	return resolver, rejector
}

func sayHelloAsync(greeting <-chan string) (<-chan string, <-chan error) {
	resolver := make(chan string, 1)
	rejector := make(chan error, 1)
	go func() {
		v := <-greeting
		time.Sleep(1 * time.Second)
		if v == "" {
			rejector <- fmt.Errorf("greeting was not set")
		} else {
			resolver <- "hello: " + 
		}
	}()
	return resolver, rejector
}

func main() {
	start := time.Now()
	fmt.Println("START", start)

	name1, nameErr1 := randomNameAsync()
	name2, nameErr2 := randomNameAsync()

	hello1, helloErr1 := sayHelloAsync(name1)
	hello2, helloErr2 := sayHelloAsync(name2)

	fmt.Println("MIDDLE", time.Since(start))

	select {
	case greetings := <-await.AllStringsAsync(hello1, hello2):
		for _, greeting := range greetings {
			fmt.Println(greeting)
		}
	case err := <-await.AnyErrorAsync(nameErr1, nameErr2, helloErr1, helloErr2):
		panic(err)
	}

	fmt.Println("END", time.Since(start))
}
```
