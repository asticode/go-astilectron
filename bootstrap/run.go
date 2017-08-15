package bootstrap

import (
	"os"
	"path/filepath"

	"github.com/asticode/go-astilectron"
	"github.com/pkg/errors"
)

// Run runs the bootstrap
func Run(o Options) (err error) {
	// Create astilectron
	var a *astilectron.Astilectron
	if a, err = astilectron.New(o.AstilectronOptions); err != nil {
		return errors.Wrap(err, "creating new astilectron failed")
	}
	defer a.Close()
	a.HandleSignals()

	// Adapt astilectron
	if o.AdaptAstilectron != nil {
		o.AdaptAstilectron(a)
	}

	// Base directory path default to executable path
	if o.BaseDirectoryPath == "" {
		if o.BaseDirectoryPath, err = os.Executable(); err != nil {
			return errors.Wrap(err, "getting executable path failed")
		}
		o.BaseDirectoryPath = filepath.Dir(o.BaseDirectoryPath)
	}

	// Provision
	if err = provision(o.BaseDirectoryPath, o.RestoreAssets, o.CustomProvision); err != nil {
		return errors.Wrap(err, "provisioning failed")
	}

	// Start
	if err = a.Start(); err != nil {
		return errors.Wrap(err, "starting astilectron failed")
	}

	// Serve or handle messages
	var url string
	if o.MessageHandler == nil {
		var ln = serve(o.BaseDirectoryPath, o.AdaptRouter, o.TemplateData)
		defer ln.Close()
		url = "http://" + ln.Addr().String() + o.Homepage
	} else {
		url = filepath.Join(o.BaseDirectoryPath, "resources", "app", o.Homepage)
	}

	// Debug
	if o.Debug {
		o.WindowOptions.Width = astilectron.PtrInt(*o.WindowOptions.Width + 700)
	}

	// Init window
	var w *astilectron.Window
	if w, err = a.NewWindow(url, o.WindowOptions); err != nil {
		return errors.Wrap(err, "new window failed")
	}

	// Adapt window
	if o.AdaptWindow != nil {
		o.AdaptWindow(w)
	}

	// Handle messages
	if o.MessageHandler != nil {
		w.On(astilectron.EventNameWindowEventMessage, handleMessages(w, o.MessageHandler))
	}

	// Create window
	if err = w.Create(); err != nil {
		return errors.Wrap(err, "creating window failed")
	}

	// Debug
	if o.Debug {
		if err = w.OpenDevTools(); err != nil {
			return errors.Wrap(err, "opening dev tools failed")
		}
	}

	// On wait
	if o.OnWait != nil {
		if err = o.OnWait(a, w); err != nil {
			return errors.Wrap(err, "onwait failed")
		}
	}

	// Blocking pattern
	a.Wait()
	return
}
