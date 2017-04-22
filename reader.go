package astilectron

import (
	"bufio"
	"bytes"
	"io"

	"encoding/json"

	"github.com/asticode/go-astilog"
)

// reader represents an object capable of reading the stdout
type reader struct {
	d *Dispatcher
	s io.ReadCloser
}

// newReader creates a new reader
func newReader(d *Dispatcher, stdout io.ReadCloser) *reader {
	return &reader{
		d: d,
		s: stdout,
	}
}

// close closes the reader properly
func (r *reader) close() error {
	return r.s.Close()
}

// read reads from stdout
func (r *reader) read() {
	var reader = bufio.NewReader(r.s)
	for {
		// Read next line
		var b []byte
		var err error
		if b, err = reader.ReadBytes('\n'); err != nil {
			if err == io.EOF {
				astilog.Debug("Electron stopped")
				r.d.Dispatch(Event{Name: EventNameElectronStopped, TargetID: mainTargetID})
				return
			} else {
				astilog.Errorf("%s while reading", err)
				continue
			}
		}
		b = bytes.TrimSpace(b)
		astilog.Debugf("Electron says: %s", b)

		// This is an astilectron message
		if bytes.HasSuffix(b, boundary) {
			// Trim boundary
			b = bytes.TrimSuffix(b, boundary)

			// Unmarshal
			var e Event
			if err = json.Unmarshal(b, &e); err != nil {
				astilog.Errorf("%s while unmarshaling %s", err, b)
				continue
			}

			// Dispatch
			r.d.Dispatch(e)
		}
	}
}
