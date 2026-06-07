package result

import (
	"errors"
	"fmt"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestResultError(t *testing.T) {
	result := Result{
		Err: errors.New("fail"),
	}

	assert.NotNil(t, result.Error())

	result = Result{
		Panic: "mock",
	}

	assert.NotNil(t, result.Error())
}

func TestResult_Error_ReturnsErr(t *testing.T) {
	result := Result{
		Err:   fmt.Errorf("fail"),
		Panic: "mock",
	}

	err := result.Error()

	assert.NotNil(t, result.Error())
	assert.Equal(t, "fail", err.Error())
}

func TestResult_Error_ReturnsPanic(t *testing.T) {
	result := Result{
		Panic: "mock",
		Stack: []byte("stacktrace"),
	}

	err := result.Error()

	assert.ErrorType[PanicError](t, result.Error())
	assert.Inside(t, "panic: mock", err.Error())
	assert.Inside(t, "stacktrace", err.Error())
}
