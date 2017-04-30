package main

import (
	"flag"
	"os"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

//go:generate go-bindata -pkg $GOPACKAGE -o vendor.go ../vendor/
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
	a.SetProvisioner(astilectron.NewDisembedderProvisioner(Asset, "../vendor/astilectron-v0.1.0.zip", "../vendor/electron-v1.6.5.zip"))
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

	// Blocking pattern
	a.Wait()
}
