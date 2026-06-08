package writer

import (
	"io"
	"os"
)

// StderrWithNewline returns an io.Writer that writes to standard error (os.Stderr) and appends a newline character after each write operation. 
// This is useful for ensuring that log messages or error outputs are properly separated when written to the console or log files.
func StderrWithNewline() io.Writer {
	return &Decorator{
		writer: os.Stderr,
		suffix: []byte("\n"),
	}
}
