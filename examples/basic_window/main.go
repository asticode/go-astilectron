package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

func main() {
	// Parse flags
	flag.Parse()

	// Set up logger
	astilog.SetLogger(astilog.New(astilog.FlagConfig()))

	// Retrieve current directory path
	var p string
	var err error
	if p, err = os.Executable(); err != nil {
		astilog.Fatal(errors.Wrap(err, "retrieving os executable failed"))
	}
	p = filepath.Dir(p)

	// Create astilectron
	var a *astilectron.Astilectron
	if a, err = astilectron.New(astilectron.Options{BaseDirectoryPath: p}); err != nil {
		astilog.Fatal(errors.Wrap(err, "creating new astilectron failed"))
	}
	defer a.Close()
	a.HandleSignals()
	a.On(astilectron.EventNameElectronStopped, func(e astilectron.Event) (deleteListener bool) {
		a.Stop()
		return
	})

	// Start
	if err = a.Start(); err != nil {
		astilog.Fatal(errors.Wrap(err, "starting failed"))
	}

	// Create window
	var w *astilectron.Window
	if w, err = a.NewWindow(&astilectron.WindowOptions{
		Center: astilectron.PtrBool(true),
		Show:   astilectron.PtrBool(false),
		Height: astilectron.PtrInt(600),
		Width:  astilectron.PtrInt(600),
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "creating new window failed"))
	}
	w.Show()

	// Blocking pattern
	a.Wait()
}
