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

	// New app menu
	// You can do the same thing with a window
	var m = a.NewMenu([]*astilectron.MenuItemOptions{
		{
			Label: astilectron.PtrStr("Separator"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Label: astilectron.PtrStr("Normal 1")},
				{Label: astilectron.PtrStr("Normal 2")},
				{Type: astilectron.MenuItemTypeSeparator},
				{Label: astilectron.PtrStr("Normal 3")},
			},
		},
		{
			Label: astilectron.PtrStr("Checkbox"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Checked: astilectron.PtrBool(true), Label: astilectron.PtrStr("Checkbox 1"), Type: astilectron.MenuItemTypeCheckbox},
				{Label: astilectron.PtrStr("Checkbox 2"), Type: astilectron.MenuItemTypeCheckbox},
				{Label: astilectron.PtrStr("Checkbox 3"), Type: astilectron.MenuItemTypeCheckbox},
			},
		},
		{
			Label: astilectron.PtrStr("Radio"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Checked: astilectron.PtrBool(true), Label: astilectron.PtrStr("Radio 1"), Type: astilectron.MenuItemTypeRadio},
				{Label: astilectron.PtrStr("Radio 2"), Type: astilectron.MenuItemTypeRadio},
				{Label: astilectron.PtrStr("Radio 3"), Type: astilectron.MenuItemTypeRadio},
			},
		},
		{
			Label: astilectron.PtrStr("Roles"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Label: astilectron.PtrStr("Minimize"), Role: astilectron.MenuItemRoleMinimize},
				{Label: astilectron.PtrStr("Close"), Role: astilectron.MenuItemRoleClose},
			},
		},
	})

	// Retrieve a menu item
	var mi *astilectron.MenuItem
	if mi, err = m.Item(1, 0); err != nil {
		astilog.Fatal(errors.Wrap(err, "fetching menu item 1,0 failed"))
	}

	// Add listener
	mi.On(astilectron.EventNameMenuItemEventClicked, func(e astilectron.Event) bool {
		astilog.Infof("Menu item has been clicked. 'Checked' status is now %t", *e.MenuItemOptions.Checked)
		return false
	})

	// Create the menu
	if err = m.Create(); err != nil {
		astilog.Fatal(errors.Wrap(err, "creating app menu failed"))
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

	// Manipulate menu item
	time.Sleep(time.Second)
	if mi, err = m.Item(1, 1); err != nil {
		astilog.Fatal(errors.Wrap(err, "fetching menu item 1,1 failed"))
	}
	if err = mi.SetChecked(true); err != nil {
		astilog.Fatal(errors.Wrap(err, "setting checked failed"))
	}

	// Insert menu item dynamically
	time.Sleep(time.Second)
	var ni = m.NewItem(&astilectron.MenuItemOptions{
		Label: astilectron.PtrStr("Inserted"),
		SubMenu: []*astilectron.MenuItemOptions{
			{Label: astilectron.PtrStr("Inserted 1")},
			{Label: astilectron.PtrStr("Inserted 2")},
		},
	})
	if err = m.Insert(1, ni); err != nil {
		astilog.Fatal(errors.Wrap(err, "inserting menu item failed"))
	}

	// Fetch sub menu
	var s *astilectron.SubMenu
	if s, err = m.SubMenu(0); err != nil {
		astilog.Fatal(errors.Wrap(err, "fetching sub menu 0 failed"))
	}

	// Append menu item dynamically
	time.Sleep(time.Second)
	ni = s.NewItem(&astilectron.MenuItemOptions{
		Label: astilectron.PtrStr("Appended"),
		SubMenu: []*astilectron.MenuItemOptions{
			{Label: astilectron.PtrStr("Appended 1")},
			{Label: astilectron.PtrStr("Appended 2")},
		},
	})
	if err = s.Append(ni); err != nil {
		astilog.Fatal(errors.Wrap(err, "appending menu item failed"))
	}

	// Pop up sub menu
	time.Sleep(time.Second)
	if err = s.Popup(&astilectron.MenuPopupOptions{PositionOptions: astilectron.PositionOptions{X: astilectron.PtrInt(50), Y: astilectron.PtrInt(50)}}); err != nil {
		astilog.Fatal(errors.Wrap(err, "popping up sub menu failed"))
	}

	// Close popup
	time.Sleep(time.Second)
	if err = s.ClosePopup(); err != nil {
		astilog.Fatal(errors.Wrap(err, "closing popup sub menu failed"))
	}

	// Blocking pattern
	a.Wait()
}
