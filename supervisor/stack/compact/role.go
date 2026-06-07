package compact

// FrameRole identifies the role of a stack frame within a processed stack trace.
type FrameRole string

const (
	// RoleOrigin marks the frame where the panic or error originated.
	RoleOrigin FrameRole = "origin"

	// RoleStack marks a regular stack frame that is part of the execution path.
	RoleStack  FrameRole = "stack"
)

// RolePriority defines the order in which frame roles should be processed or displayed.
var RolePriority = []FrameRole{
	RoleOrigin,
	RoleStack,
}
