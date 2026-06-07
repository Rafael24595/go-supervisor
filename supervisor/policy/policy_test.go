package policy

import (
	"testing"
	"time"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-supervisor/supervisor/result"
)

func TestDefaultPolicyContract(t *testing.T) {
	policy := Default()

	assert.Equal(t, time.Second, policy.Delay)
	assert.Equal(t, 3, policy.MaxRestarts)
	assert.NotNil(t, policy.RestartIf)
	assert.True(t, policy.RestartIf(result.Result{}))
}

func TestRestartAllowed(t *testing.T) {
	policy := Policy{
		MaxRestarts: 3,
	}

	tests := []struct {
		name     string
		restarts uint
		want     bool
	}{
		{"no restarts", 0, true},
		{"few restarts", 2, true},
		{"limit restarts", 3, false},
		{"too many restarts", 9, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := policy.RestartAllowed(tt.restarts)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUnlimitedRestarts(t *testing.T) {
	p := Policy{
		MaxRestarts: UnlimitedRestarts,
	}

	for i := range 1000 {
		assert.True(t, p.RestartAllowed(uint(i)))
	}
}

func TestRestartIf(t *testing.T) {
	called := uint8(0)

	p := Policy{
		RestartIf: func(res result.Result) bool {
			called += 1
			return false
		},
	}

	assert.False(t, p.RestartIf(result.Result{}))
	assert.Equal(t, 1, called)
}
