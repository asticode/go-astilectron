package astilectron

import (
	"net/http"
	"os"
	"path/filepath"

	"context"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/context"
	"github.com/asticode/go-astitools/http"
	"github.com/asticode/go-astitools/zip"
	"github.com/pkg/errors"
)

// Download is a cancellable function that downloads a src into a dst using a specific *http.Client and deals with
// failed downloads
func Download(ctx context.Context, c *http.Client, dst, src string) (err error) {
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

// Unzip unzips a src into a dst.
// Possible src formats are /path/to/zip.zip or /path/to/zip.zip/internal/path.
func Unzip(ctx context.Context, dst, src string) error {
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

// synchronousFunc executes a function and blocks until it has received a specific event or the canceller has been
// cancelled
func synchronousFunc(c *asticontext.Canceller, l listenable, fn func(), eventNameDone string) {
	var ctx, cancel = c.NewContext()
	defer cancel()
	l.On(eventNameDone, func(e Event) (deleteListener bool) {
		cancel()
		return true
	})
	fn()
	<-ctx.Done()
}

// synchronousEvent sends an event and blocks until it has received a specific event or the canceller has been cancelled
func synchronousEvent(c *asticontext.Canceller, l listenable, w *writer, e Event, eventNameDone string) (err error) {
	synchronousFunc(c, l, func() {
		if err = w.write(e); err != nil {
			err = errors.Wrapf(err, "writing %+v event failed", e)
			return
		}
		return
	}, eventNameDone)
	return
}
