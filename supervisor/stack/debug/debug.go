package debug

import (
	"runtime/debug"

	"github.com/Rafael24595/go-supervisor/supervisor/stack"
)

// Stack returns a stack provider that captures the current stack trace using the runtime/debug package.
func Stack() stack.Provider {
	return func() []byte {
		return debug.Stack()
	}
}
