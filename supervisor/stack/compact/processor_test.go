package compact

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestClearFunction(t *testing.T) {
	got := clearFunction(
		"github.com/acme/project",
		"github.com/acme/project/pkg/foo.Bar",
	)

	assert.Equal(t, "pkg/foo", got)
}

func TestClearFunctionRootPackage(t *testing.T) {
	got := clearFunction(
		"github.com/acme/project",
		"github.com/acme/project.Main",
	)

	assert.Equal(t, "", got)
}

func TestClearFunctionWithoutMethod(t *testing.T) {
	got := clearFunction(
		"github.com/acme/project",
		"github.com/acme/project/pkg/foo",
	)

	assert.Equal(t, "pkg/foo", got)
}

func TestClearFunctionMethod(t *testing.T) {
	got := clearFunction(
		"github.com/acme/project",
		"github.com/acme/project/pkg/foo.(*Service).Run",
	)

	assert.Equal(t, "pkg/foo", got)
}

func TestProcessWithoutAbsoluteRoot(t *testing.T) {
	p := standardProcessor{
		hasPath: false,
	}

	got := p.Process("/tmp/file.go")

	assert.Equal(t, "/tmp/file.go", got)
}

func TestProcessRelativePath(t *testing.T) {
	p := standardProcessor{
		path:    "/home/user/project/",
		hasPath: true,
	}

	got := p.Process(
		"/home/user/project/pkg/foo/file.go",
	)

	assert.Equal(t, "pkg/foo/file.go", got)
}

func TestProcessExternalPath(t *testing.T) {
	p := standardProcessor{
		path:    "/home/user/project/",
		hasPath: true,
	}

	got := p.Process(
		"/usr/local/go/src/runtime/panic.go",
	)

	assert.Equal(
		t,
		"/usr/local/go/src/runtime/panic.go",
		got,
	)
}
