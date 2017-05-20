package bootstrap

import (
	"os"
	"path/filepath"

	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// provision provisions the resources as well as the custom provision
func provision(fnA RestoreAssets, fnP CustomProvision) (pr string, err error) {
	// Get executable path
	var p string
	if p, err = os.Executable(); err != nil {
		err = errors.Wrap(err, "getting executable path failed")
		return
	}
	p = filepath.Dir(p)

	// Provision resources
	pr = filepath.Join(p, "resources")
	if _, err = os.Stat(pr); os.IsNotExist(err) {
		// Restore assets
		astilog.Debugf("Restoring assets in %s", p)
		if err = fnA(p, "resources"); err != nil {
			err = errors.Wrapf(err, "restoring assets in %s failed", p)
			return
		}
	} else if err != nil {
		err = errors.Wrapf(err, "stating %s failed", pr)
		return
	}

	// Custom provision
	if fnP != nil {
		if err = fnP(); err != nil {
			err = errors.Wrap(err, "custom provisioning failed")
			return
		}
	}
	return
}
