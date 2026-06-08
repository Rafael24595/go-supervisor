package writer

import (
	"bytes"
	"sync"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestDecorator_Success(t *testing.T) {
	var baseBuf bytes.Buffer

	decorator := Decorator{
		writer: &baseBuf,
		prefix: []byte("[PRE] "),
		suffix: []byte(" [SUF]"),
	}

	input := []byte("HELLO")
	n, err := decorator.Write(input)

	assert.Nil(t, err)
	assert.Size(t, n, input)
	assert.Equal(t, "[PRE] HELLO [SUF]", baseBuf.String())
}

func TestDecorator_AdjustWritten(t *testing.T) {
	tests := []struct {
		name       string
		prefixLen  int
		expected   int
		written    int
		wantResult int
	}{
		{
			name:       "Fail before prefix ends",
			prefixLen:  5,
			expected:   10,
			written:    2,
			wantResult: 0,
		},
		{
			name:       "Fail just after prefix ends",
			prefixLen:  5,
			expected:   10,
			written:    5,
			wantResult: 0,
		},
		{
			name:       "Fail in the middle of the payload",
			prefixLen:  5,
			expected:   10,
			written:    9,
			wantResult: 4,
		},
		{
			name:       "Fail at the end of the payload",
			prefixLen:  5,
			expected:   10,
			written:    15,
			wantResult: 10,
		},
		{
			name:       "Fail in the middle of the suffix",
			prefixLen:  5,
			expected:   10,
			written:    17,
			wantResult: 10,
		},
	}

	d := &Decorator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d.prefix = make([]byte, tt.prefixLen)
			got := d.adjustWritten(tt.written, tt.expected)

			assert.Equal(t, tt.wantResult, got)
		})
	}
}

func TestWriterDecorator_Concurrency(t *testing.T) {
	var baseBuf bytes.Buffer

	decorator := Decorator{
		writer: &baseBuf,
	}

	var wg sync.WaitGroup

	goroutines := 50
	writes := 20
	message := []byte("A")

	for range goroutines {
		wg.Go(func() {
			for range writes {
				_, _ = decorator.Write(message)
			}
		})
	}

	wg.Wait()

	expectedBytes := goroutines * writes * len(message)
	assert.Equal(t, expectedBytes, baseBuf.Len())
}
