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

	// Start loader
	if o.StartLoader != nil {
		o.StartLoader(a)
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

	// Serve
	var ln = serve(o.BaseDirectoryPath, o.AdaptRouter, o.TemplateData)
	defer ln.Close()

	// Debug
	if o.Debug {
		o.WindowOptions.Width = astilectron.PtrInt(*o.WindowOptions.Width + 700)
	}

	// Init window
	var w *astilectron.Window
	if w, err = a.NewWindow("http://"+ln.Addr().String()+o.Homepage, o.WindowOptions); err != nil {
		return errors.Wrap(err, "new window failed")
	}

	// Adapt window
	if o.AdaptWindow != nil {
		o.AdaptWindow(w)
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

	// Blocking pattern
	a.Wait()
	return
}
