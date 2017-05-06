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

	// Get base dir path
	var err error
	var p = os.Getenv("GOPATH") + "/src/github.com/asticode/go-astilectron/examples"

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

	// Create window in the primary display
	var w *astilectron.Window
	if w, err = a.NewWindowInDisplay(p+"/index.html", &astilectron.WindowOptions{
		Icon:   astilectron.PtrStr(os.Getenv("GOPATH") + "/src/github.com/asticode/go-astilectron/examples/6.icons/gopher.png"),
		Height: astilectron.PtrInt(600),
		Show:   astilectron.PtrBool(false),
		Width:  astilectron.PtrInt(600),
	}, a.PrimaryDisplay()); err != nil {
		astilog.Fatal(errors.Wrap(err, "new window failed"))
	}
	if err = w.Create(); err != nil {
		astilog.Fatal(errors.Wrap(err, "creating window failed"))
	}
	if err = w.Center(); err != nil {
		astilog.Fatal(errors.Wrap(err, "centering window failed"))
	}
	if err = w.Show(); err != nil {
		astilog.Fatal(errors.Wrap(err, "showing window failed"))
	}

	// Move window to the second display if any
	var displays = a.Displays()
	var display = displays[0]
	if len(displays) > 1 {
		display = displays[1]
	}
	time.Sleep(time.Second)
	if err = w.MoveInDisplay(display, 50, 50); err != nil {
		astilog.Fatal(errors.Wrap(err, "moving window in display failed"))
	}

	// Blocking pattern
	a.Wait()
}
