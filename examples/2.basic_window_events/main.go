package main

import (
	"flag"
	"os"
	"time"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

func main() {
	// Parse flags
	flag.Parse()

	// Set up logger
	astilog.SetLogger(astilog.New(astilog.FlagConfig()))

	// Create astilectron
	var a *astilectron.Astilectron
	var err error
	if a, err = astilectron.New(astilectron.Options{BaseDirectoryPath: os.Getenv("GOPATH") + "/src/github.com/asticode/go-astilectron/examples"}); err != nil {
		astilog.Fatal(errors.Wrap(err, "creating new astilectron failed"))
	}
	defer a.Close()
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		astilog.Fatal(errors.Wrap(err, "starting failed"))
	}

	// Create window
	var w *astilectron.Window
	if w, err = a.NewWindow("http://google.com", &astilectron.WindowOptions{
		Center: astilectron.PtrBool(true),
		Height: astilectron.PtrInt(600),
		Width:  astilectron.PtrInt(600),
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "new window failed"))
	}
	if err = w.Create(); err != nil {
		astilog.Fatal(errors.Wrap(err, "creating window failed"))
	}

	// Add listener
	w.On(astilectron.EventNameWindowEventMove, func(e astilectron.Event) (deleteListener bool) {
		astilog.Info("Window moved")
		return
	})

	// Simulate actions
	go func() {
		time.Sleep(time.Second)
		if err = w.Move(0, 0); err != nil {
			astilog.Fatal(errors.Wrap(err, "moving window failed"))
		}
		time.Sleep(time.Second)
		if err = w.Resize(200, 200); err != nil {
			astilog.Fatal(errors.Wrap(err, "resizing window failed"))
		}
		time.Sleep(time.Second)
		if err = w.Maximize(); err != nil {
			astilog.Fatal(errors.Wrap(err, "maximizing window failed"))
		}
		time.Sleep(time.Second)
		if err = w.Unmaximize(); err != nil {
			astilog.Fatal(errors.Wrap(err, "unmaximizing window failed"))
		}
		time.Sleep(time.Second)
		if err = w.Minimize(); err != nil {
			astilog.Fatal(errors.Wrap(err, "minimizing window failed"))
		}
		time.Sleep(time.Second)
		if err = w.Restore(); err != nil {
			astilog.Fatal(errors.Wrap(err, "restoring window failed"))
		}
		time.Sleep(time.Second)
		if err = w.Resize(600, 600); err != nil {
			astilog.Fatal(errors.Wrap(err, "resizing window failed"))
		}
		if err = w.Center(); err != nil {
			astilog.Fatal(errors.Wrap(err, "centering window failed"))
		}
	}()

	// Blocking pattern
	a.Wait()
}
