package astilectron

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/context"
	"github.com/asticode/go-astitools/http"
	"github.com/asticode/go-astitools/io"
	"github.com/asticode/go-astitools/zip"
	"github.com/pkg/errors"
)

// Download is a cancellable function that downloads a src into a dst using a specific *http.Client and deals with
// failed downloads
func Download(ctx context.Context, c *http.Client, src, dst string) (err error) {
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

	// Download
	if err = astihttp.Download(ctx, c, src, dst); err != nil {
		return errors.Wrap(err, "astihttp.Download failed")
	}

	// We need to close the file manually before removing it
	fp.Close()

	// Remove dst.processing file
	if err = os.Remove(dstProcessing); err != nil {
		return errors.Wrapf(err, "removing %s failed", dstProcessing)
	}
	return
}

// Disembed is a cancellable disembed of an src to a dst using a custom Disembedder
func Disembed(ctx context.Context, d Disembedder, src, dst string) (err error) {
	// Log
	astilog.Debugf("Disembedding %s into %s...", src, dst)

	// No need to disembed
	if _, err = os.Stat(dst); err != nil && !os.IsNotExist(err) {
		return errors.Wrapf(err, "stating %s failed", dst)
	} else if err == nil {
		astilog.Debugf("%s already exists, skipping disembed...", dst)
		return
	}
	err = nil

	// Make sure directory exists
	var dirPath = filepath.Dir(dst)
	astilog.Debugf("Creating %s", dirPath)
	if err = os.MkdirAll(dirPath, 0755); err != nil {
		return errors.Wrapf(err, "mkdirall %s failed", dirPath)
	}

	// Create dst
	var f *os.File
	astilog.Debugf("Creating %s", dst)
	if f, err = os.Create(dst); err != nil {
		err = errors.Wrapf(err, "creating %s failed", dst)
		return
	}
	defer f.Close()

	// Disembed
	var b []byte
	astilog.Debugf("Disembedding %s", src)
	if b, err = d(src); err != nil {
		err = errors.Wrapf(err, "disembedding %s failed", src)
		return
	}

	// Copy
	astilog.Debugf("Copying disembedded data to %s", dst)
	if _, err = astiio.Copy(ctx, bytes.NewReader(b), f); err != nil {
		err = errors.Wrapf(err, "copying disembedded data into %s failed", dst)
		return
	}
	return
}

// Unzip unzips a src into a dst.
// Possible src formats are /path/to/zip.zip or /path/to/zip.zip/internal/path.
func Unzip(ctx context.Context, src, dst string) error {
	astilog.Debugf("Unzipping %s into %s", src, dst)
	return astizip.Unzip(ctx, src, dst)
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

// synchronousFunc executes a function, blocks until it has received a specific event or the canceller has been
// cancelled and returns the corresponding event
func synchronousFunc(c *asticontext.Canceller, l listenable, fn func(), eventNameDone string) (e Event) {
	var ctx, cancel = c.NewContext()
	defer cancel()
	l.On(eventNameDone, func(i Event) (deleteListener bool) {
		e = i
		cancel()
		return true
	})
	fn()
	<-ctx.Done()
	return
}

// synchronousEvent sends an event, blocks until it has received a specific event or the canceller has been cancelled
// and returns the corresponding event
func synchronousEvent(c *asticontext.Canceller, l listenable, w *writer, i Event, eventNameDone string) (o Event, err error) {
	o = synchronousFunc(c, l, func() {
		if err = w.write(i); err != nil {
			err = errors.Wrapf(err, "writing %+v event failed", i)
			return
		}
		return
	}, eventNameDone)
	return
}
