package compact

import (
	"runtime"
	"strings"

	"github.com/Rafael24595/go-supervisor/supervisor/stack"
)

const defaultCaller = "supervisor.runProtected.func1"

// FramesRoles groups stack frames by their assigned role.
type FramesRoles map[FrameRole][]stackFrame

type stackFrame struct {
	Function string
	File     string
	Line     uint
}

// Stack returns a stack provider that captures the current stack trace and formats it using the provided formatter and options.
func Stack(formatter Formatter, opts ...option) stack.Provider {
	cfg := newConfig(formatter, opts...)

	return func() []byte {
		return compactProvider(cfg)
	}
}

func compactProvider(cfg config) []byte {
	frames := getStackFrame(cfg.caller)
	return cfg.formatter(frames, cfg.processor)
}

func getStackFrame(caller string) FramesRoles {
	frames := getRuntimeFrames()

	var frame runtime.Frame
	hasMore := true

	result := make(FramesRoles)

	firstFrame := true
	callerFound := false

	for hasMore {
		frame, hasMore = frames.Next()
		if strings.Contains(frame.File, "runtime/") {
			continue
		}

		if !callerFound {
			callerFound = strings.HasSuffix(frame.Function, caller)
			continue
		}

		role := RoleStack
		if firstFrame {
			role = RoleOrigin
		}

		result[role] = append(result[role], stackFrame{
			Function: frame.Function,
			File:     frame.File,
			Line:     uint(frame.Line),
		})

		firstFrame = false
	}

	return result
}

func getRuntimeFrames() *runtime.Frames {
	counters := make([]uintptr, 64)
	entries := runtime.Callers(2, counters)
	return runtime.CallersFrames(counters[:entries])
}
