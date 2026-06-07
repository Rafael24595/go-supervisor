package result

// Result represents the outcome of a supervised task execution, including any error or panic that occurred.
type Result struct {
	// Err holds any error returned by the task. It is nil if the task completed successfully or panicked.
	Err error

	// Panic holds the value passed to panic if the task panicked. It is nil if the task completed successfully or returned an error.
	Panic any

	// Stack holds the stack trace at the point of panic. It is nil if the task completed successfully or returned an error.
	Stack []byte
}

// Error returns an error representation of the Result.
//
// If there was an error, it returns that error.
// If there was a panic, it returns an error containing the panic message and stack trace.
// If there were no errors or panics, it returns nil.
func (r Result) Error() error {
	if r.Err != nil {
		return r.Err
	}

	if r.Panic != nil {
		return PanicError{
			Panic: r.Panic,
			Stack: r.Stack,
		}
	}

	return nil
}
