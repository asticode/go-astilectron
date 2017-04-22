package astilectron

import (
	"io"

	"encoding/json"

	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// writer represents an object capable of writing in the stdin
type writer struct {
	s io.WriteCloser
}

// newWriter creates a new writer
func newWriter(stdin io.WriteCloser) *writer {
	return &writer{
		s: stdin,
	}
}

// close closes the writer properly
func (r *writer) close() error {
	return r.s.Close()
}

// write writes to the stdin
func (r *writer) write(e Event) (err error) {
	// Marshal
	var b []byte
	if b, err = json.Marshal(e); err != nil {
		return errors.Wrapf(err, "Marshaling %+v failed", e)
	}

	// Write
	var m = append(b, boundary...)
	m = append(m, '\n')
	astilog.Debugf("Sending to Electron: %s", m)
	if _, err = r.s.Write(m); err != nil {
		return errors.Wrapf(err, "Writing %s failed", string(m))
	}
	return
}
