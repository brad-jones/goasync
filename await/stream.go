package await

import (
	"github.com/brad-jones/goasync/v2/task"
)

// Stream will wait for the first task to return and continue to do that until
// all tasks have returned.
//
// For example:
// 	s := await.Stream(foo(), bar())
// 	for s.Wait() {
// 		r, err := s.Result()
// 	}
func Stream(awaitables ...*task.Task) *StreamInstance {
	return &StreamInstance{awaitables: awaitables}
}

// StreamInstance is the object that is returned by Stream
type StreamInstance struct {
	awaitables []*task.Task
	current    *task.Task
}

// Wait will return true once a task has finished & then remove that task from
// the list of tasks to wait for. It will return false when there are no more
// tasks to wait for.
func (s *StreamInstance) Wait() bool {
	if len(s.awaitables) == 0 {
		return false
	}

	doneCh := make(chan struct{}, 1)
	awaitableCh := make(chan *task.Task, 1)
	for _, v := range s.awaitables {
		awaitable := v
		go func() {
			select {
			case <-*awaitable.Done:
				awaitableCh <- awaitable
				close(awaitableCh)
			case <-doneCh:
				return
			}
		}()
	}
	s.current = <-awaitableCh
	close(doneCh)

	newAwaitables := []*task.Task{}
	for _, v := range s.awaitables {
		if v != s.current {
			newAwaitables = append(newAwaitables, v)
		}
	}
	s.awaitables = newAwaitables

	return true
}

// Result is an alias for the completed task's Result method.
func (s *StreamInstance) Result() (interface{}, error) {
	return s.current.Result()
}

// MustResult is an alias for the completed task's MustResult method.
func (s *StreamInstance) MustResult() interface{} {
	return s.current.MustResult()
}

// Task return the completed task
func (s *StreamInstance) Task() *task.Task {
	return s.current
}
