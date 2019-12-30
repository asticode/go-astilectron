package astilectron

import (
	"bytes"
	"context"
	"os"
	"path/filepath"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Download is a cancellable function that downloads a src into a dst using a specific *http.Client and cleans up on
// failed downloads
func Download(ctx context.Context, d *astikit.HTTPDownloader, src, dst string) (err error) {
	// Log
	astilog.Debugf("Downloading %s into %s", src, dst)

	// Destination already exists
	if _, err = os.Stat(dst); err == nil {
		astilog.Debugf("%s already exists, skipping download...", dst)
		return
	} else if !os.IsNotExist(err) {
		return errors.Wrapf(err, "stating %s failed", dst)
	}
	err = nil

	// Clean up on error
	defer func(err *error) {
		if *err != nil || ctx.Err() != nil {
			astilog.Debugf("Removing %s...", dst)
			os.Remove(dst)
		}
	}(&err)

	// Make sure the dst directory  exists
	if err = os.MkdirAll(filepath.Dir(dst), 0775); err != nil {
		return errors.Wrapf(err, "mkdirall %s failed", filepath.Dir(dst))
	}

	// Download
	if err = d.DownloadInFile(ctx, dst, astikit.HTTPDownloaderSrc{URL: src}); err != nil {
		return errors.Wrap(err, "DownloadInFile failed")
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

	// Clean up on error
	defer func(err *error) {
		if *err != nil || ctx.Err() != nil {
			astilog.Debugf("Removing %s...", dst)
			os.Remove(dst)
		}
	}(&err)

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
		return errors.Wrapf(err, "creating %s failed", dst)
	}
	defer f.Close()

	// Disembed
	var b []byte
	astilog.Debugf("Disembedding %s", src)
	if b, err = d(src); err != nil {
		return errors.Wrapf(err, "disembedding %s failed", src)
	}

	// Copy
	astilog.Debugf("Copying disembedded data to %s", dst)
	if _, err = astikit.Copy(ctx, f, bytes.NewReader(b)); err != nil {
		return errors.Wrapf(err, "copying disembedded data into %s failed", dst)
	}
	return
}

// Unzip unzips a src into a dst.
// Possible src formats are /path/to/zip.zip or /path/to/zip.zip/internal/path.
func Unzip(ctx context.Context, src, dst string) (err error) {
	// Clean up on error
	defer func(err *error) {
		if *err != nil || ctx.Err() != nil {
			astilog.Debugf("Removing %s...", dst)
			os.RemoveAll(dst)
		}
	}(&err)

	// Unzipping
	astilog.Debugf("Unzipping %s into %s", src, dst)
	if err = astikit.Unzip(ctx, dst, src); err != nil {
		err = errors.Wrapf(err, "unzipping %s into %s failed", src, dst)
		return
	}
	return
}

// synchronousFunc executes a function, blocks until it has received a specific event or the canceller has been
// cancelled and returns the corresponding event
func synchronousFunc(parentCtx context.Context, l listenable, fn func(), eventNameDone string) (e Event) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	l.On(eventNameDone, func(i Event) (deleteListener bool) {
		if ctx.Err() == nil {
			e = i
		}
		cancel()
		return true
	})
	if fn != nil {
		fn()
	}
	<-ctx.Done()
	return
}

// synchronousEvent sends an event, blocks until it has received a specific event or the canceller has been cancelled
// and returns the corresponding event
func synchronousEvent(ctx context.Context, l listenable, w *writer, i Event, eventNameDone string) (o Event, err error) {
	o = synchronousFunc(ctx, l, func() {
		if err = w.write(i); err != nil {
			err = errors.Wrapf(err, "writing %+v event failed", i)
			return
		}
	}, eventNameDone)
	return
}
