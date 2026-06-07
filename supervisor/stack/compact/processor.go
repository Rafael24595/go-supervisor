package compact

import (
	"runtime"
	"runtime/debug"
	"strings"
)

// Processor transforms a stack frame path into a different representation.
type Processor func(path string) string

type standardProcessor struct {
	absolute    string
	hasAbsolute bool
}

// EmptyProcessor returns a Processor that leaves paths unchanged.
func EmptyProcessor() Processor {
	instance := standardProcessor{
		hasAbsolute: false,
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
	if !ok {
		return EmptyProcessor()
	}

	pc, file, _, ok := runtime.Caller(0)
	if !ok {
		return EmptyProcessor()
	}

	packageRelative := clearFunction(
		info.Main.Path,
		runtime.FuncForPC(pc).Name(),
	)

	absolute, _, ok := strings.Cut(file, packageRelative)
	if !ok {
		return EmptyProcessor()
	}

	instance := standardProcessor{
		absolute:    absolute,
		hasAbsolute: true,
	}

	return instance.Process
}

func clearFunction(project, function string) string {
	functionRelative := strings.ReplaceAll(function, project, "")
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
	if !r.hasAbsolute {
		return path
	}

	_, relative, ok := strings.Cut(path, r.absolute)
	if !ok {
		return path
	}

	return relative
}
