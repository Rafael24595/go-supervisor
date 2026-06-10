package compact

import (
	"runtime"
	"runtime/debug"
	"strings"
)

// Processor transforms a stack frame path into a different representation.
type Processor func(path string) string

type standardProcessor struct {
	path    string
	hasPath bool
}

// EmptyProcessor returns a Processor that leaves paths unchanged.
func EmptyProcessor() Processor {
	instance := standardProcessor{
		hasPath: false,
	}
	return instance.Process
}

// StandardProcessor returns a Processor that converts absolute file paths
// into paths relative to the current project when possible.
//
// If the project root cannot be determined, the returned Processor leaves
// paths unchanged.
func StandardProcessor() Processor {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Path == "" {
		return EmptyProcessor()
	}

	frame, ok := findMainFrame(info.Main.Path)
	if !ok {
		return EmptyProcessor()
	}

	mainDir := clearFunction(
		info.Main.Path,
		frame.Function,
	)

	path, _, ok := strings.Cut(frame.File, mainDir)
	if !ok {
		return EmptyProcessor()
	}

	instance := standardProcessor{
		path:    strings.ToLower(path),
		hasPath: true,
	}

	return instance.Process
}

func findMainFrame(main string) (runtime.Frame, bool) {
	size := 64

	for {
		pcs := make([]uintptr, size)
		n := runtime.Callers(2, pcs)

		frames := runtime.CallersFrames(pcs[:n])

		for {
			frame, more := frames.Next()
			if strings.HasPrefix(frame.Function, main) {
				return frame, true
			}

			if !more {
				break
			}
		}

		if n < size {
			return runtime.Frame{}, false
		}

		size *= 2
	}
}

func clearFunction(project, function string) string {
	functionRelative := strings.TrimPrefix(function, project)
	functionRelative = strings.TrimPrefix(functionRelative, "/")

	lastSlash := strings.LastIndex(functionRelative, "/")
	if lastSlash == -1 {
		lastSlash = 0
	}

	dotAfterSlash := strings.Index(functionRelative[lastSlash:], ".")
	if dotAfterSlash == -1 {
		return functionRelative
	}

	return functionRelative[:lastSlash+dotAfterSlash]
}

// Process converts an absolute path into a project-relative path when the
// processor has been initialized with a project root. Otherwise, it returns
// the original path unchanged.
func (r standardProcessor) Process(path string) string {
	if !r.hasPath {
		return path
	}

	_, relative, ok := strings.Cut(
		strings.ToLower(path), r.path,
	)

	if !ok {
		return path
	}

	return relative
}
