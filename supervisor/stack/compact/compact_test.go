package compact

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

const caller = "compact.Stack.func1"

func TestStack_ReturnsSomething(t *testing.T) {
	provider := Stack(func(frames FramesRoles, p Processor) []byte {
		return []byte("ok")
	})

	out := provider()

	assert.Equal(t, "ok", string(out))
}

func TestGetStackFrame_FiltersRuntimeFrames(t *testing.T) {
	frames := getStackFrame("compact.getStackFrame")

	for _, list := range frames {
		for _, f := range list {
			assert.False(t, strings.Contains(f.File, "runtime/"))
		}
	}
}

func TestGetStackFrame_HasRoles(t *testing.T) {
	frames := getStackFrame("compact.getStackFrame")

	_, hasOrigin := frames[RoleOrigin]
	_, hasStack := frames[RoleStack]

	assert.True(t, hasOrigin || hasStack)
}

func TestRolePriority(t *testing.T) {
	frames := getStackFrame(caller)

	for i := 0; i < len(frames[RoleOrigin]); i++ {
		assert.Equal(t, RoleOrigin, RolePriority[0])
	}
}

func TestStack_EndToEnd(t *testing.T) {
	provider := Stack(
		GroupedFormatter,
		withCaller(caller),
	)

	out := provider()

	lines := strings.Split(string(out), "\n")
	for _, v := range lines {
		println(v)
	}

	assert.True(t, len(out) > 0)
	assert.Equal(t, "ORIGIN", lines[0])
	assert.True(t, strings.HasPrefix(lines[1], "	compact.TestStack_EndToEnd"))

	assert.Equal(t, "", lines[2])

	assert.Equal(t, "STACK", lines[3])
	assert.True(t, strings.HasPrefix(lines[4], "	testing.tRunner"))
}
