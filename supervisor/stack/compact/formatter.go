package compact

import (
	"fmt"
	"strings"
)

// Formatter formats a collection of stack frames into a compact textual
// representation.
type Formatter func(FramesRoles, Processor) []byte

// GroupedFormatter formats frames grouped by their role.
//
// Each role is rendered as a separate section headed by the role name,
// followed by the frames that belong to that role.
//
// Example:
//
//	ORIGIN
//		main.main [cmd/app/main.go:10]
//
//	STACK
//		service.Run [internal/service/run.go:42]
//		handler.Execute [internal/handler/execute.go:17]
func GroupedFormatter(frames FramesRoles, processor Processor) []byte {
	compact := make([]string, 0, len(RolePriority))

	for _, role := range RolePriority {
		frames, ok := frames[role]
		if !ok {
			continue
		}

		roleStack := make([]string, 0, len(frames)+1)

		roleStack = append(roleStack, strings.ToUpper(string(role)))

		for _, frame := range frames {
			fragments := strings.Split(frame.Function, "/")
			function := fragments[len(fragments)-1]

			file := processor(frame.File)

			roleStack = append(roleStack,
				fmt.Sprintf("	%s [%s:%d]", function, file, frame.Line),
			)
		}

		compact = append(compact,
			strings.Join(roleStack, "\n"),
		)
	}

	return []byte(
		strings.Join(compact, "\n\n"),
	)
}

// InlineFormatter formats frames as a flat list.
//
// Frames marked as RoleOrigin are prefixed with "->" while all other
// frames use a neutral marker.
//
// Example:
//
//	-> main.main [cmd/app/main.go:10]
//	   service.Run [internal/service/run.go:42]
//	   handler.Execute [internal/handler/execute.go:17]
func InlineFormatter(frames FramesRoles, processor Processor) []byte {
	compact := make([]string, 0, len(frames))
	for _, role := range RolePriority {
		frames, ok := frames[role]
		if !ok {
			continue
		}

		for _, frame := range frames {
			marker := "   "
			if role == RoleOrigin {
				marker = "-> "
			}

			function := processFunction(frame.Function)
			file := processor(frame.File)

			formatted := fmt.Sprintf("%s%s [%s:%d]", marker, function, file, frame.Line)
			compact = append(compact, formatted)
		}
	}

	return []byte(
		strings.Join(compact, "\n"),
	)
}

func processFunction(function string) string {
	lastSlash := strings.LastIndex(function, "/")
	if lastSlash == -1 {
		return function
	}

	return strings.TrimPrefix(function[lastSlash:], "/")
}
