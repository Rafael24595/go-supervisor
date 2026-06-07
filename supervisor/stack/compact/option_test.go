package compact

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestDefaultConfig(t *testing.T) {
	mockPath := "github.com/acme/project/pkg/foo.(*Service).Run"

	formatter := func(FramesRoles, Processor) []byte {
		return []byte("golang")
	}

	cfg := defaultConfig(formatter)

	assert.Equal(t, defaultCaller, cfg.caller)

	assert.NotNil(t, cfg.caller)
	assert.Equal(t,
		string(formatter(FramesRoles{}, cfg.processor)),
		string(cfg.formatter(FramesRoles{}, cfg.processor)),
	)

	assert.NotNil(t, cfg.processor)
	assert.Equal(t,
		string(EmptyProcessor()(mockPath)),
		string(cfg.processor(mockPath)),
	)
}

func TestWithProcessor(t *testing.T) {
	cfg := defaultConfig(nil)

	expected := func(path string) string {
		return "mock:" + path
	}

	opt := WithProcessor(expected)
	opt(&cfg)

	got := cfg.processor("test.go")
	assert.Equal(t, "mock:test.go", got)
}
