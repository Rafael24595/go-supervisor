package debug

import (
	"github.com/Rafael24595/go-supervisor/supervisor/stack"
)

// Stack returns a stack provider that returns an empty stack trace.
func Stack() stack.Provider {
	return func() []byte {
		return make([]byte, 0)
	}
}
