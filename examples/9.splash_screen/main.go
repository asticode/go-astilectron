package main

import (
	"flag"
	"os"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astisplash"
	"github.com/pkg/errors"
)

func main() {
	// Parse flags
	flag.Parse()

	// Set up logger
	astilog.SetLogger(astilog.New(astilog.FlagConfig()))

	// Get base dir path
	var p = os.Getenv("GOPATH") + "/src/github.com/asticode/go-astilectron/examples"

	// Build splasher
	var s *astisplash.Splasher
	var err error
	if s, err = astisplash.New(); err != nil {
		astilog.Fatal(errors.Wrap(err, "building splasher failed"))
	}
	defer s.Close()

	// Splash
	var sp *astisplash.Splash
	if sp, err = s.Splash(p+"/splash.png", 400, 400); err != nil {
		astilog.Fatal(errors.Wrap(err, "splashing failed"))
	}

	// Create astilectron
	var a *astilectron.Astilectron
	if a, err = astilectron.New(astilectron.Options{
		AppName:            "Astilectron",
		AppIconDefaultPath: p + "/gopher.png",
		AppIconDarwinPath:  p + "/gopher.icns",
		BaseDirectoryPath:  p,
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "creating new astilectron failed"))
	}
	defer a.Close()
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		astilog.Fatal(errors.Wrap(err, "starting failed"))
	}

	// Close splash
	if err = sp.Close(); err != nil {
		astilog.Fatal(errors.Wrap(err, "closing splash failed"))
	}

	// Create window
	var w *astilectron.Window
	if w, err = a.NewWindow(p+"/index.html", &astilectron.WindowOptions{
		Center: astilectron.PtrBool(true),
		Height: astilectron.PtrInt(600),
		Width:  astilectron.PtrInt(600),
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "new window failed"))
	}
	if err = w.Create(); err != nil {
		astilog.Fatal(errors.Wrap(err, "creating window failed"))
	}

	// Blocking pattern
	a.Wait()
}
