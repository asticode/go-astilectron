package astilectron

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockedWriter represents a mocked writer
type mockedWriter struct {
	c  bool
	fn func()
	w  []string
	wg *sync.WaitGroup
}

// Close implements the io.Closer interface
func (w *mockedWriter) Close() error {
	w.c = true
	return nil
}

// Write implements io.Writer interface
func (w *mockedWriter) Write(p []byte) (int, error) {
	w.w = append(w.w, string(p))
	if w.fn != nil {
		w.fn()
	}
	if w.wg != nil {
		w.wg.Done()
	}
	return len(p), nil
}

// TestWriter tests the writer
func TestWriter(t *testing.T) {
	// Init
	var mw = &mockedWriter{}
	var w = newWriter(mw, &logger{})

	// Test write
	err := w.write(Event{Name: "test", TargetID: "target_id"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"{\"name\":\"test\",\"targetID\":\"target_id\"}\n"}, mw.w)

	// Test close
	err = w.close()
	assert.NoError(t, err)
	assert.True(t, mw.c)
}
