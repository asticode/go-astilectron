package astilectron

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Download downloads a src into a dst using a specific *http.Client
func Download(c *http.Client, dst, src string) (err error) {
	// Log
	astilog.Debugf("Downloading %s into %s", src, dst)

	// Make sure the directory of the dst exists
	if err = os.MkdirAll(filepath.Dir(dst), 0775); err != nil {
		return errors.Wrapf(err, "mkdirall %s failed", filepath.Dir(dst))
	}

	// Check whether dst and dst.processing exist
	var dstProcessing = dst + ".processing"
	var dstExists, dstProcessingExists = true, true
	if _, err = os.Stat(dst); os.IsNotExist(err) {
		dstExists = false
	} else if err != nil {
		return errors.Wrapf(err, "stating %s failed", dst)
	}
	if _, err = os.Stat(dstProcessing); os.IsNotExist(err) {
		dstProcessingExists = false
	} else if err != nil {
		return errors.Wrapf(err, "stating %s failed", dstProcessing)
	}
	err = nil

	// Skipping download
	if dstExists && !dstProcessingExists {
		astilog.Debugf("%s already exists, skipping download...", dst)
		return
	} else if dstProcessingExists {
		astilog.Debugf("%s already exists, cleaning up and downloading again...", dstProcessing)
		for _, p := range []string{dst, dstProcessing} {
			if err = os.Remove(p); err != nil {
				return errors.Wrapf(err, "removing %s failed", p)
			}
		}
	}

	// Create the dst.processing file
	var fp *os.File
	if fp, err = os.Create(dstProcessing); err != nil {
		return errors.Wrapf(err, "creating file %s failed", dstProcessing)
	}
	defer fp.Close()

	// Create the dst file
	var f *os.File
	if f, err = os.Create(dst); err != nil {
		return errors.Wrapf(err, "creating file %s failed", dst)
	}
	defer f.Close()

	// Send request
	var resp *http.Response
	if resp, err = c.Get(src); err != nil {
		return errors.Wrapf(err, "getting %s failed", src)
	}
	defer resp.Body.Close()

	// Validate status code
	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(err, "getting %s returned %d status code", src, resp.StatusCode)
	}

	// Copy
	if _, err = io.Copy(f, resp.Body); err != nil {
		return errors.Wrapf(err, "copying content from %s to %s failed", src, dst)
	}

	// Remove dst.processing file
	if err = os.Remove(dstProcessing); err != nil {
		return errors.Wrapf(err, "removing %s failed", dstProcessing)
	}
	return
}

// Unzip unzips a src into a dst
func Unzip(dst, src string) (err error) {
	// Log
	astilog.Debugf("Unzipping %s into %s", src, dst)

	// Open overall reader
	var r *zip.ReadCloser
	if r, err = zip.OpenReader(src); err != nil {
		return errors.Wrapf(err, "opening overall zip reader on %s failed", src)
	}
	defer r.Close()

	// Loop through files
	for _, f := range r.File {
		// Open file reader
		var fr io.ReadCloser
		if fr, err = f.Open(); err != nil {
			return errors.Wrapf(err, "opening zip reader on file %s failed", f.Name)
		}
		defer fr.Close()

		// Only unzip files
		var p = filepath.Join(dst, f.Name)
		if !f.FileInfo().IsDir() {
			// Make sure the directory of the file exists
			if err = os.MkdirAll(filepath.Dir(p), 0775); err != nil {
				return errors.Wrapf(err, "mkdirall %s failed", filepath.Dir(p))
			}

			// Open the file
			var fl *os.File
			if fl, err = os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode()); err != nil {
				return errors.Wrapf(err, "opening file %s failed", p)
			}
			defer fl.Close()

			// Copy
			if _, err = io.Copy(fl, fr); err != nil {
				return errors.Wrapf(err, "copying %s into %s failed", f.Name, p)
			}
		}
	}
	return
}

// PtrBool transforms a bool into a *bool
func PtrBool(i bool) *bool {
	return &i
}

// PtrInt transforms an int into an *int
func PtrInt(i int) *int {
	return &i
}

// PtrStr transforms a string into a *string
func PtrStr(i string) *string {
	return &i
}

// synchronousFunc executes a function and blocks until it has received a specific event
func synchronousFunc(l Listenable, eventNameDone string, fn func()) {
	var c = make(chan bool)
	defer func() {
		if c != nil {
			close(c)
		}
	}()
	l.On(eventNameDone, func(e Event) (deleteListener bool) {
		close(c)
		c = nil
		return true
	})
	fn()
	<-c
}

// synchronousEvent sends an event and blocks until it has received a specific event
func synchronousEvent(l Listenable, w *writer, e Event, eventNameDone string) (err error) {
	synchronousFunc(l, eventNameDone, func() {
		if err = w.write(e); err != nil {
			err = errors.Wrapf(err, "writing %+v event failed", e)
			return
		}
		return
	})
	return
}
