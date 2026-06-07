package result

import (
	"errors"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestPanicErrorErrorIncludesPanic(t *testing.T) {
	err := PanicError{
		Panic: "mock",
	}

	assert.Inside(t, "mock", err.Error())
}

func TestPanicErrorErrorIncludesStack(t *testing.T) {
	err := PanicError{
		Panic: "mock",
		Stack: []byte("fake-stack"),
	}

	assert.Inside(t, "fake-stack", err.Error())
}

func TestPanicErrorUnwrapReturnsNil(t *testing.T) {
	err := PanicError{
		Panic: "mock",
	}

	assert.Nil(t, err.Unwrap())
}

func TestPanicErrorUnwrapReturnsOriginalError(t *testing.T) {
	expected := errors.New("mock")

	err := PanicError{
		Panic: expected,
	}

	assert.ErrorIs(t, expected, err)
}

func TestPanicErrorErrorsAs(t *testing.T) {
	err := PanicError{
		Panic: errors.New("mock"),
	}

	assert.ErrorType[PanicError](t, err)
}

type customError struct{}

func (customError) Error() string {
	return "custom"
}

func TestPanicErrorUnwrapTypedError(t *testing.T) {
	expected := customError{}

	err := PanicError{
		Panic: expected,
	}

	assert.ErrorIs(t, expected, err)
}

func TestPanicErrorWithEmptyStack(t *testing.T) {
	err := PanicError{
		Panic: "mock",
	}

	assert.NotEqual(t, "", err.Error())
}
