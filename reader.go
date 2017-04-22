package astilectron

import (
	"bytes"
	"io"

	"bufio"

	"github.com/asticode/go-astilog"
)

// Vars
var (
	boundary = []byte("+++astilectron_boundary+++")
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
				r.d.Dispatch(Event{Name: EventNameElectronStop, TargetID: mainTargetID})
				return
			} else {
				astilog.Errorf("%s while reading", err)
			}
		}

		astilog.Debug(string(b))

		// This is an astilectron message
		if bytes.HasSuffix(b, boundary) {
			// TODO
		} else {
			r.d.Dispatch(Event{Name: EventNameElectronLog, Payload: string(b), TargetID: mainTargetID})
		}
	}
}
