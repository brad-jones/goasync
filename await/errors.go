package await

// ErrTaskFailed is returned by the All methods when at least one task returns an error.
type ErrTaskFailed struct {
	Errors []error
}

func (e *ErrTaskFailed) Error() string {
	return "await: one or more errors were returned from the awaited tasks"
}
