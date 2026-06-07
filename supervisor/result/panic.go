package result

import "fmt"

// PanicError represents an error that occurred due to a panic in a supervised task, including the panic value and stack trace.
type PanicError struct {
	// Panic holds the value passed to panic. It can be of any type, but is often an error or a string.
	Panic any

	// Stack holds the stack trace at the point of panic, which can be useful for debugging.
	Stack []byte
}

// Error returns a string representation of the PanicError, including the panic value and stack trace.
func (e PanicError) Error() string {
	return fmt.Sprintf("panic: %v\n%s", e.Panic, e.Stack)
}

// Unwrap allows PanicError to be treated as an error for compatibility with errors.Is and errors.As.
func (e PanicError) Unwrap() error {
	if err, ok := e.Panic.(error); ok {
		return err
	}

	return nil
}
