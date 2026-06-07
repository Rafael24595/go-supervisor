package supervisor

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-supervisor/supervisor/policy"
	"github.com/Rafael24595/go-supervisor/supervisor/result"
)

func TestDefaultConfig(t *testing.T) {
	name := "mock"
	cfg := defaultConfig(name)

	assert.Equal(t, name, cfg.name)
	assert.Equal(t, context.Background(), cfg.ctx)
	assert.Equal[io.Writer](t, os.Stderr, cfg.writer)

	policy := policy.Default()
	assert.Equal(t, policy.Delay, cfg.policy.Delay)
	assert.Equal(t, policy.MaxRestarts, cfg.policy.MaxRestarts)

	assert.True(t, cfg.stack != nil)

	rawStack := cfg.stack()
	assert.True(t, len(rawStack) > 0)
	assert.True(t, strings.Contains(string(rawStack), "TestDefaultConfig"))
}

func TestWithContext(t *testing.T) {
	cfg := defaultConfig("mock")

	ctx := t.Context()

	opt := WithContext(ctx)
	opt(&cfg)

	assert.Equal(t, ctx, cfg.ctx)
}

func TestWithPolicy(t *testing.T) {
	cfg := defaultConfig("mock")

	delay := time.Second * 2
	maxRestarts := uint(10)

	res := result.Result{
		Err: errors.New("mocked"),
	}

	opt := WithPolicy(policy.Policy{
		Delay:       delay,
		MaxRestarts: maxRestarts,
		RestartIf: func(res result.Result) bool {
			return res.Err.Error() == "mocked"
		},
	})
	opt(&cfg)

	assert.Equal(t, delay, cfg.policy.Delay)
	assert.Equal(t, maxRestarts, cfg.policy.MaxRestarts)
	assert.True(t, cfg.policy.RestartIf(res))
}

func TestWithWriter(t *testing.T) {
	cfg := defaultConfig("mock")

	writer := io.Discard

	opt := WithWriter(writer)
	opt(&cfg)

	assert.Equal(t, writer, cfg.writer)
}

func TestNewConfig_Integration(t *testing.T) {
	ctx := t.Context()
	writer := io.Discard
	policy := policy.Policy{MaxRestarts: 5}

	cfg := newConfig("full-mock",
		WithContext(ctx),
		WithWriter(writer),
		WithPolicy(policy))

	assert.Equal(t, "full-mock", cfg.name)
	assert.Equal(t, ctx, cfg.ctx)
	assert.Equal(t, writer, cfg.writer)
	assert.Equal(t, policy.MaxRestarts, cfg.policy.MaxRestarts)
}

func TestWithStackProvider(t *testing.T) {
	cfg := defaultConfig("mock")

	expectedStack := "mocked stack trace data"
	mockProvider := func() []byte {
		return []byte(expectedStack)
	}

	opt := WithStackProvider(mockProvider)
	opt(&cfg)

	assert.Equal(t, expectedStack, string(cfg.stack()))
}
