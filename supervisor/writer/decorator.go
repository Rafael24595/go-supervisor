package writer

import (
	"bytes"
	"io"
	"sync"
)

// Decorator is an io.Writer that decorates the output with a prefix and suffix. It is designed to be thread-safe, allowing concurrent writes without data races. 
// The Decorator can be used to add consistent formatting to log messages or other outputs, ensuring that each write operation is properly framed with the specified prefix and suffix.
type Decorator struct {
	mu     sync.Mutex
	writer io.Writer
	prefix []byte
	suffix []byte
	buffer bytes.Buffer
}

// Write implements the io.Writer interface for the Decorator. It writes the provided byte slice to the underlying writer, while adding the specified prefix and suffix.
// The method ensures that the write operation is thread-safe by using a mutex to synchronize access to the internal buffer and the underlying writer. 
// It returns the number of bytes from the input slice that were successfully written, along with any error encountered during the write operation.
func (s *Decorator) Write(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.buffer.Reset()

	s.buffer.Write(s.prefix)
	s.buffer.Write(p)
	s.buffer.Write(s.suffix)

	written, err := s.writer.Write(s.buffer.Bytes())
	if err == nil {
		return len(p), nil
	}

	n = s.adjustWritten(written, len(p))
	return n, err
}

func (s *Decorator) adjustWritten(written, expected int) int {
	n := written - len(s.prefix)
	if n < 0 {
		return 0
	}

	if n > expected {
		return expected
	}

	return n
}
