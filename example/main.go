package main

import (
	"flag"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	astiptr "github.com/asticode/go-astitools/ptr"
	"github.com/pkg/errors"
)

func main() {
	// Parse flags
	astilog.SetHandyFlags()
	flag.Parse()
	astilog.FlagInit()

	// Create astilectron
	a, err := astilectron.New(astilectron.Options{AppName: "Test"})
	if err != nil {
		astilog.Fatal(errors.Wrap(err, "main: creating astilectron failed"))
	}
	defer a.Close()

	// Handle signals
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		astilog.Fatal(errors.Wrap(err, "main: starting astilectron failed"))
	}

	// New window
	var w *astilectron.Window
	if w, err = a.NewWindow("example/index.html", &astilectron.WindowOptions{
		Center:          astiptr.Bool(true),
		Height:          astiptr.Int(700),
		Width:           astiptr.Int(700),
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "main: new window failed"))
	}

	// Create windows
	if err = w.Create(); err != nil {
		astilog.Fatal(errors.Wrap(err, "main: creating window failed"))
	}

	// Blocking pattern
	a.Wait()
}
