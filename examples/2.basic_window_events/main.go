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
	a.On(astilectron.EventNameAppClose, func(e astilectron.Event) (deleteListener bool) {
		a.Stop()
		return
	})

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
	w.Create()

	// Add listener
	w.On(astilectron.EventNameWindowEventMove, func(e astilectron.Event) (deleteListener bool) {
		astilog.Info("Window moved")
		return
	})

	// Simulate actions
	go func() {
		time.Sleep(time.Second)
		w.Move(0, 0)
		time.Sleep(time.Second)
		w.Resize(200, 200)
		time.Sleep(time.Second)
		w.Maximize()
		time.Sleep(time.Second)
		w.Unmaximize()
		time.Sleep(time.Second)
		w.Minimize()
		time.Sleep(time.Second)
		w.Restore()
		time.Sleep(time.Second)
		w.Resize(600, 600)
		w.Center()
	}()

	// Blocking pattern
	a.Wait()
}
