# goasync

[![PkgGoDev](https://pkg.go.dev/badge/github.com/brad-jones/goasync/v2)](https://pkg.go.dev/github.com/brad-jones/goasync/v2)
[![GoReport](https://goreportcard.com/badge/github.com/brad-jones/goasync/v2)](https://goreportcard.com/report/github.com/brad-jones/goasync/v2)
[![GoLang](https://img.shields.io/badge/golang-%3E%3D%201.15.1-lightblue.svg)](https://golang.org)
![.github/workflows/main.yml](https://github.com/brad-jones/goasync/workflows/.github/workflows/main.yml/badge.svg?branch=v2)
[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![KeepAChangelog](https://img.shields.io/badge/Keep%20A%20Changelog-1.0.0-%23E05735)](https://keepachangelog.com/)
[![License](https://img.shields.io/github/license/brad-jones/goasync.svg)](https://github.com/brad-jones/goerr/blob/v2/LICENSE)

Package goasync is a helper framework for writing asynchronous code in go.
It's primary goal is to reduce the amount of boiler plate code one has to
write to do concurrent tasks.

_Looking for v1, see the [master branch](https://github.com/brad-jones/goasync/tree/master)_

## Quick Start

`go get -u github.com/brad-jones/goasync/v2/...`

```go
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
```

Running the above will output something similar to:

```
START 2020-09-11 18:44:14.2406928 +1000 AEST m=+0.003027901
doing work
doing work
END 1.0013651s
```

_Also see further working examples under: <https://github.com/brad-jones/goasync/tree/v2/examples>_
