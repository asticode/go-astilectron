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
	var ln = serve(o.BaseDirectoryPath, o.TemplateData)
	defer ln.Close()

	// Create window
	var w *astilectron.Window
	if w, err = a.NewWindow("http://"+ln.Addr().String()+o.Homepage, o.WindowOptions); err != nil {
		return errors.Wrap(err, "new window failed")
	}
	if err = w.Create(); err != nil {
		return errors.Wrap(err, "creating window failed")
	}

	// Blocking pattern
	a.Wait()
	return
}
