package stack

// Package stack provides an abstraction for capturing and formatting stack traces in Go. 
// It defines a Provider type, which is a function that returns a byte slice representing the stack trace. 
// This allows for flexibility in how stack traces are captured and formatted, 
// enabling users to implement custom stack providers or use existing ones, such as the debug package's Stack function.
type Provider func() []byte
