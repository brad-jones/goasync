package main

import (
	"fmt"
	"time"

	"github.com/brad-jones/goasync/v2/task"
	uuid "github.com/satori/go.uuid"
)

func main() {
	start := time.Now()
	fmt.Println("START", start)

	task.New(func(t *task.Internal) {
		time.Sleep(1 * time.Second)
		t.Resolve(uuid.NewV4().String())
	}).Then(func(result interface{}, t *task.Internal) {
		time.Sleep(1 * time.Second)
		fmt.Println("hello: " + result.(string))
	}).MustWait()

	fmt.Println("END", time.Since(start))
}
