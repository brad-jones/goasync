package task

// ErrStoppingTaskTimeout is returned by StopWithTimeout when the given duration 
// has passed and the task has still not stopped.
type ErrStoppingTaskTimeout struct {
}

func (e *ErrStoppingTaskTimeout) Error() string {
	return "task: stopping task took too long to stop"
}