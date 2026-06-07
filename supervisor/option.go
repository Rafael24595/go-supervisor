package supervisor

import (
	"context"
	"io"
	"os"

	"github.com/Rafael24595/go-supervisor/supervisor/policy"
	"github.com/Rafael24595/go-supervisor/supervisor/stack"
	"github.com/Rafael24595/go-supervisor/supervisor/stack/debug"
)

type option func(*config)

type config struct {
	name   string
	ctx    context.Context
	policy policy.Policy
	writer io.Writer
	stack  stack.Provider
}

func newConfig(name string, options ...option) config {
	cfg := defaultConfig(name)
	for _, o := range options {
		o(&cfg)
	}
	return cfg
}

func defaultConfig(name string) config {
	return config{
		name:   name,
		ctx:    context.Background(),
		policy: policy.Default(),
		writer: os.Stderr,
		stack:  debug.Stack(),
	}
}

// WithContext allows setting a custom context for the supervisor, which can be used to control cancellation and timeouts for the supervised tasks.
func WithContext(ctx context.Context) option {
	return func(c *config) {
		c.ctx = ctx
	}
}

// WithPolicy allows setting a custom policy for the supervisor, which can be used to control the behavior of the supervised tasks.
func WithPolicy(policy policy.Policy) option {
	return func(c *config) {
		c.policy = policy
	}
}

// WithWriter allows setting a custom writer for the supervisor, which can be used to redirect logs and output from the supervised tasks.
func WithWriter(writer io.Writer) option {
	return func(c *config) {
		c.writer = writer
	}
}

// WithStackProvider allows setting a custom stack provider for the supervisor, which can be used to customize how stack traces are captured and formatted when a task panics.
func WithStackProvider(stack stack.Provider) option {
	return func(c *config) {
		c.stack = stack
	}
}
