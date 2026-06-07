package policy

import (
	"time"

	"github.com/Rafael24595/go-supervisor/supervisor/result"
)

// UnlimitedRestarts is a constant that can be used to indicate that there is no limit on the number of restarts.
const UnlimitedRestarts = uint(0)

// Policy defines the rules for restarting a supervised task.
type RestartIfFunc func(res result.Result) bool

// Always is a RestartIfFunc that always returns true, indicating that the task should always be restarted regardless of the result.
func Always(result.Result) bool { return true }

// Never is a RestartIfFunc that always returns false, indicating that the task should never be restarted regardless of the result.
func Never(result.Result) bool { return false }

// Policy defines the rules for restarting a supervised task.
type Policy struct {
	// Delay is the duration to wait before restarting a failed task.
	Delay time.Duration

	// MaxRestarts is the maximum number of times a task can be restarted before giving up. A value of 0 means unlimited restarts.
	MaxRestarts uint

	// RestartIf is a function that determines whether a task should be restarted based on the result of its execution.
	RestartIf RestartIfFunc
}

// New creates a new Policy with the specified parameters.
func New(
	delay time.Duration,
	maxRestarts uint,
	restartIf RestartIfFunc,
) Policy {
	return Policy{
		Delay:       delay,
		MaxRestarts: maxRestarts,
		RestartIf:   restartIf,
	}
}

// Default returns a Policy with default settings: a delay of 1 second, a maximum of 3 restarts, and a RestartIf function that always returns true.
func Default() Policy {
	return Policy{
		Delay:       time.Second,
		MaxRestarts: 3,
		RestartIf:   Always,
	}
}

// RestartAllowed checks if the number of restarts is within the allowed limit defined by MaxRestarts.
func (p Policy) RestartAllowed(restarts uint) bool {
	if p.MaxRestarts == 0 {
		return true
	}
	return restarts < p.MaxRestarts
}
