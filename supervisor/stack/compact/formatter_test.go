package compact

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestGroupedFormatter_SingleRoleSingleFrame(t *testing.T) {
	frames := FramesRoles{
		RoleOrigin: {
			{
				Function: "main.main",
				File:     "/project/main.go",
				Line:     10,
			},
		},
	}

	processor := func(p string) string {
		return "main.go"
	}

	out := GroupedFormatter(frames, processor)
	lines := strings.Split(string(out), "\n")

	assert.Equal(t, "ORIGIN", lines[0])
	assert.Equal(t, "\tmain.main [main.go:10]", lines[1])
}

func TestGroupedFormatter_MultipleFrames(t *testing.T) {
	frames := FramesRoles{
		RoleStack: {
			{Function: "a.b", File: "f1.go", Line: 1},
			{Function: "c.d", File: "f2.go", Line: 2},
		},
	}

	processor := func(p string) string { return p }

	out := GroupedFormatter(frames, processor)
	lines := strings.Split(string(out), "\n")

	assert.Equal(t, "STACK", lines[0])
	assert.Equal(t, "\ta.b [f1.go:1]", lines[1])
	assert.Equal(t, "\tc.d [f2.go:2]", lines[2])
}

func TestGroupedFormatter_RolePriorityOrder(t *testing.T) {
	frames := FramesRoles{
		RoleStack: {
			{Function: "stack.fn", File: "s.go", Line: 1},
		},
		RoleOrigin: {
			{Function: "origin.fn", File: "o.go", Line: 1},
		},
	}

	processor := func(p string) string { return p }

	out := GroupedFormatter(frames, processor)
	lines := strings.Split(string(out), "\n")

	assert.Equal(t, "ORIGIN", lines[0])
	assert.Equal(t, "\torigin.fn [o.go:1]", lines[1])

	assert.Equal(t, "", lines[2])

	assert.Equal(t, "STACK", lines[3])
	assert.Equal(t, "\tstack.fn [s.go:1]", lines[4])
}

func TestInlineFormatterr_SingleRoleSingleFrame(t *testing.T) {
	frames := FramesRoles{
		RoleOrigin: {
			{
				Function: "main.main",
				File:     "/project/main.go",
				Line:     10,
			},
		},
	}

	processor := func(p string) string {
		return "main.go"
	}

	out := InlineFormatter(frames, processor)
	lines := strings.Split(string(out), "\n")

	assert.Equal(t, "-> main.main [main.go:10]", lines[0])
}

func TestInlineFormatter_MultipleFrames(t *testing.T) {
	frames := FramesRoles{
		RoleStack: {
			{Function: "a.b", File: "f1.go", Line: 1},
			{Function: "c.d", File: "f2.go", Line: 2},
		},
	}

	processor := func(p string) string { return p }

	out := InlineFormatter(frames, processor)
	lines := strings.Split(string(out), "\n")

	assert.Equal(t, "   a.b [f1.go:1]", lines[0])
	assert.Equal(t, "   c.d [f2.go:2]", lines[1])
}

func TestInlineFormatter_RolePriorityOrder(t *testing.T) {
	frames := FramesRoles{
		RoleStack: {
			{Function: "stack.fn", File: "s.go", Line: 1},
		},
		RoleOrigin: {
			{Function: "origin.fn", File: "o.go", Line: 1},
		},
	}

	processor := func(p string) string { return p }

	out := InlineFormatter(frames, processor)
	lines := strings.Split(string(out), "\n")

	assert.Equal(t, "-> origin.fn [o.go:1]", lines[0])
	assert.Equal(t, "   stack.fn [s.go:1]", lines[1])
}

func TestFormatter_FunctionParsing(t *testing.T) {
	frames := FramesRoles{
		RoleStack: {
			{Function: "github.com/pkg/service.Run", File: "f.go", Line: 1},
		},
	}

	out := InlineFormatter(frames, func(s string) string { return s })

	assert.Equal(t, "   service.Run [f.go:1]", string(out))
}

func TestFormatter_UsesProcessor(t *testing.T) {
	called := false

	frames := FramesRoles{
		RoleStack: {
			{Function: "a.b", File: "file.go", Line: 1},
		},
	}

	processor := func(p string) string {
		called = true
		return "X"
	}

	GroupedFormatter(frames, processor)

	assert.True(t, called)
}

func TestFormatter_EmptyInput(t *testing.T) {
	out := GroupedFormatter(FramesRoles{}, func(s string) string { return s })

	assert.Equal(t, "", string(out))
}