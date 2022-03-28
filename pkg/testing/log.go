package testing

import (
	"fmt"
	"testing"
)

type StdLogger struct {
	*testing.T
}

func (t StdLogger) Print(v ...interface{}) {
	fmt.Print(v)
}

func (t StdLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v)
}

// AdaptTestLogger transforms an testing.T into a StdLogger if needed
func AdaptTestLogger(t *testing.T) StdLogger {
	return StdLogger{t}
}
