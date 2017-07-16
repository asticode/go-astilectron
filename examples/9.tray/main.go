package main

import (
	"flag"
	"os"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

func main() {
	// Parse flags
	flag.Parse()

	// Set up logger
	astilog.FlagInit()

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

	// New tray
	var t = a.NewTray(&astilectron.TrayOptions{
		Image:   astilectron.PtrStr(p + "/gopher.png"),
		Tooltip: astilectron.PtrStr("Tray's tooltip"),
	})

	// New tray menu
	var m = t.NewMenu([]*astilectron.MenuItemOptions{
		{
			Label: astilectron.PtrStr("Root 1"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Label: astilectron.PtrStr("Item 1")},
				{Label: astilectron.PtrStr("Item 2")},
				{Type: astilectron.MenuItemTypeSeparator},
				{Label: astilectron.PtrStr("Item 3")},
			},
		},
		{
			Label: astilectron.PtrStr("Root 2"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Label: astilectron.PtrStr("Item 1")},
				{Label: astilectron.PtrStr("Item 2")},
			},
		},
	})

	// Create the menu
	if err = m.Create(); err != nil {
		astilog.Fatal(errors.Wrap(err, "creating tray menu failed"))
	}

	// Create tray
	if err = t.Create(); err != nil {
		astilog.Fatal(errors.Wrap(err, "creating tray failed"))
	}

	// Blocking pattern
	a.Wait()
}
