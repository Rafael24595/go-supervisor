package supervisor

import (
	"fmt"
	"time"

	"github.com/Rafael24595/go-supervisor/supervisor/result"
)

// Task represents a unit of work that can be supervised.
// It is defined as a function that returns an error, which indicates the success or failure of the task.
type Task func() error

// Supervisor manages the execution of a Task according to a specified Policy.
// It handles task restarts based on the outcome of the task and the rules defined in the Policy.
type Supervisor struct {
	name string
	opts []option
}

// New creates a new Supervisor with the given name and options.
// The name is used for logging purposes, while the options allow customization of the supervisor's behavior,
// such as setting a custom context, policy, writer, or stack provider.
func New(name string, opts ...option) *Supervisor {
	return &Supervisor{
		name: name,
		opts: opts,
	}
}

// Run executes the provided Task under the supervision of the Supervisor.
// It will handle task restarts based on the outcome of the task and the rules defined in the Policy.
// The method will continue to restart the task according to the policy until the context is cancelled or the restart limit is exceeded.
func (s *Supervisor) Run(task Task) error {
	cfg := newConfig(s.name, s.opts...)

	restarts := uint(0)

	for {
		res := runProtected(cfg, task)

		switch {
		case cfg.ctx.Err() != nil:
			fmt.Fprintf(cfg.writer, "[%s] stopping: context cancelled", cfg.name)
			return nil

		case res.Panic != nil:
			fmt.Fprintf(cfg.writer, "[%s] panic recovered: %v\n%s", cfg.name, res.Panic, res.Stack)

		case res.Err != nil:
			fmt.Fprintf(cfg.writer, "[%s] task failed: %v", cfg.name, res.Err)

		default:
			fmt.Fprintf(cfg.writer, "[%s] task exited normally", cfg.name)
			return nil
		}

		if !cfg.policy.RestartIf(res) {
			fmt.Fprintf(cfg.writer, "[%s] stop by policy", cfg.name)
			return res.Error()
		}

		restarts += 1
		if !cfg.policy.RestartAllowed(restarts) {
			fmt.Fprintf(cfg.writer, "[%s] restart limit exceeded (%d), stopping", cfg.name, restarts)
			return res.Error()
		}

		fmt.Fprintf(cfg.writer, "[%s] restarting (%d) in %s", cfg.name, restarts, cfg.policy.Delay)

		select {
		case <-cfg.ctx.Done():
			fmt.Fprintf(cfg.writer, "[%s] stopping: context cancelled", cfg.name)
			return nil

		case <-time.After(cfg.policy.Delay):
		}
	}
}

func runProtected(cfg config, task Task) (res result.Result) {
	defer func() {
		if r := recover(); r != nil {
			res.Panic = r
			res.Stack = cfg.stack()
		}
	}()

	res.Err = task()

	return
}
