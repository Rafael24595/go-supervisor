package supervisor

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-supervisor/supervisor/policy"
	"github.com/Rafael24595/go-supervisor/supervisor/result"
)

func TestRunProtectedSuccess(t *testing.T) {
	res := runProtected(
		defaultConfig("mock"),
		func() error {
			return nil
		},
	)

	assert.Nil(t, res.Error())
}

func TestRunProtectedError(t *testing.T) {
	expected := errors.New("mock")

	res := runProtected(
		defaultConfig("mock"),
		func() error {
			return expected
		},
	)

	assert.ErrorIs(t, expected, res.Err)
	assert.Nil(t, res.Panic)
}

func TestRunProtectedPanic(t *testing.T) {
	expected := errors.New("mock")

	res := runProtected(
		defaultConfig("mock"),
		func() error {
			panic(expected)
		},
	)

	assert.ErrorType[result.PanicError](t, res.Error())
	assert.ErrorIs(t, expected, res.Error())
}

func TestSupervisorRunSuccess(t *testing.T) {
	err := New("test").Run(
		func() error {
			return nil
		},
	)

	assert.Nil(t, err)
}

func TestSupervisorRunStopByPolicy(t *testing.T) {
	expected := errors.New("mock")

	sup := New(
		"test",
		WithPolicy(policy.Policy{
			RestartIf: func(result.Result) bool {
				return false
			},
		}),
		WithWriter(io.Discard),
	)

	err := sup.Run(func() error {
		return expected
	})

	assert.ErrorIs(t, expected, err)
}

func TestSupervisorRunRestarts(t *testing.T) {
	calls := 0

	sup := New(
		"test",
		WithWriter(io.Discard),
	)

	err := sup.Run(
		func() error {
			calls++
			if calls < 3 {
				return errors.New("mock")
			}

			return nil
		},
	)

	assert.Nil(t, err)
	assert.Equal(t, 3, calls)
}

func TestSupervisorRunRestartLimit(t *testing.T) {
	expected := errors.New("mock")
	calls := 0

	sup := New(
		"test",
		WithPolicy(policy.Policy{
			Delay:       0,
			MaxRestarts: 2,
			RestartIf: func(result.Result) bool {
				return true
			},
		}),
		WithWriter(io.Discard),
	)

	err := sup.Run(
		func() error {
			calls++
			return expected
		},
	)

	if err == nil {
		t.Fatal()
	}

	assert.ErrorIs(t, expected, err)
	assert.Equal(t, 2, calls)
}

func TestSupervisorRunRecoverPanic(t *testing.T) {
	calls := 0

	sup := New(
		"test",
		WithWriter(io.Discard),
	)

	err := sup.Run(func() error {
		calls++

		if calls == 1 {
			panic("mock")
		}

		return nil
	})

	assert.Nil(t, err)
	assert.Equal(t, 2, calls)
}

func TestSupervisorStopsWhenContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	sup := New(
		"test",
		WithContext(ctx),
		WithPolicy(policy.Policy{
			Delay:       time.Hour,
			MaxRestarts: policy.UnlimitedRestarts,
			RestartIf:   policy.Always,
		}),
		WithWriter(io.Discard),
	)

	calls := 0

	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	err := sup.Run(func() error {
		calls++
		return errors.New("fail")
	})

	assert.Nil(t, err)
	assert.GreaterThan(t, 0, calls)
}
