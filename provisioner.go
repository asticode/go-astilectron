package astilectron

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/io"
	"github.com/pkg/errors"
)

// Provisioner represents an object capable of provisioning Astilectron
type Provisioner interface {
	Provision(ctx context.Context, p Paths) error
}

// Default provisioner
var DefaultProvisioner = &defaultProvisioner{
	httpClient: &http.Client{},
}

// defaultProvisioner represents the default provisioner
type defaultProvisioner struct {
	httpClient *http.Client // We need to set up our own client in case we need to tweak some options such as timeout or proxy
}

// Provision implements the provisioner interface
func (p *defaultProvisioner) Provision(ctx context.Context, paths Paths) (err error) {
	// Provision astilectron
	if err = p.provisionAstilectron(ctx, paths); err != nil {
		err = errors.Wrap(err, "default provisioning astilectron failed")
		return
	}

	// Provision electron
	if err = p.provisionElectron(ctx, paths); err != nil {
		err = errors.Wrap(err, "default provisioning electron failed")
		return
	}
	return
}

// provisionAstilectron provisions astilectron
func (p *defaultProvisioner) provisionAstilectron(ctx context.Context, paths Paths) error {
	return p.provisionDownloadableZipFile(ctx, "Astilectron", paths.AstilectronApplication(), paths.AstilectronDownloadSrc(), paths.AstilectronDownloadDst(), paths.AstilectronUnzipSrc(), paths.AstilectronDirectory())
}

// provisionElectron provisions electron
func (p *defaultProvisioner) provisionElectron(ctx context.Context, paths Paths) error {
	return p.provisionDownloadableZipFile(ctx, "Electron", paths.ElectronExecutable(), paths.ElectronDownloadSrc(), paths.ElectronDownloadDst(), paths.ElectronUnzipSrc(), paths.ElectronDirectory())
}

// provisionDownloadableZipFile provisions a downloadable .zip file
func (p *defaultProvisioner) provisionDownloadableZipFile(ctx context.Context, name, pathExists, pathDownloadSrc, pathDownloadDst, pathUnzipSrc, pathDirectory string) (err error) {
	// Log
	astilog.Debugf("Default provisioning %s...", name)

	// We need to provision
	if _, err = os.Stat(pathExists); os.IsNotExist(err) {
		// Download the .zip file
		if err = Download(ctx, p.httpClient, pathDownloadDst, pathDownloadSrc); err != nil {
			return errors.Wrapf(err, "downloading %s into %s failed", pathDownloadSrc, pathDownloadDst)
		}

		// Remove previous install
		astilog.Debugf("Removing %s", pathDirectory)
		if err = os.RemoveAll(pathDirectory); err != nil && !os.IsNotExist(err) {
			return errors.Wrapf(err, "removing %s failed", pathDirectory)
		}

		// Create directory
		astilog.Debugf("Creating %s", pathDirectory)
		if err = os.MkdirAll(pathDirectory, 0755); err != nil {
			return errors.Wrapf(err, "mkdirall %s failed", pathDirectory)
		}

		// Unzip
		if err = Unzip(ctx, pathDirectory, pathUnzipSrc); err != nil {
			return errors.Wrapf(err, "unzipping %s into %s failed", pathUnzipSrc, pathDirectory)
		}
	} else if err != nil {
		return errors.Wrapf(err, "stating %s failed", pathExists)
	} else {
		astilog.Debugf("%s already exists, skipping %s default provision...", pathExists, name)
	}
	return
}

// Disembedder is a functions that allows to disembed data from a path
type Disembedder func(src string) ([]byte, error)

// NewDisembedderProvisioner creates a provisioner that can provision based on embedded data
func NewDisembedderProvisioner(d Disembedder, pathAstilectron, pathElectron string) Provisioner {
	return &disembedderProvisioner{d: d, pathAstilectron: pathAstilectron, pathElectron: pathElectron}
}

// disembedderProvisioner represents the disembedder provisioner
type disembedderProvisioner struct {
	d                             Disembedder
	pathAstilectron, pathElectron string
}

// Provision implements the provisioner interface
func (p *disembedderProvisioner) Provision(ctx context.Context, paths Paths) (err error) {
	// Disembed astilectron
	if err = p.disembedAstilectron(ctx, paths); err != nil {
		err = errors.Wrap(err, "disembedding astilectron failed")
		return
	}

	// Disembed electron
	if err = p.disembedElectron(ctx, paths); err != nil {
		err = errors.Wrap(err, "disembedding electron failed")
		return
	}

	// Default provisioner
	return DefaultProvisioner.Provision(ctx, paths)
}

// disembedAstilectron provisions astilectron
func (p *disembedderProvisioner) disembedAstilectron(ctx context.Context, paths Paths) error {
	return p.disembed(ctx, "Astilectron", p.pathAstilectron, paths.AstilectronDownloadDst())
}

// provisionElectron provisions electron
func (p *disembedderProvisioner) disembedElectron(ctx context.Context, paths Paths) error {
	return p.disembed(ctx, "Electron", p.pathElectron, paths.ElectronDownloadDst())
}

// disembed disembeds data from a src to a dst
func (p *disembedderProvisioner) disembed(ctx context.Context, name, src, dst string) (err error) {
	// Log
	astilog.Debugf("Disembedding %s...", name)

	// We need to disembed
	if _, err = os.Stat(dst); os.IsNotExist(err) {
		// Create directory
		var dirPath = filepath.Dir(dst)
		astilog.Debugf("Creating %s", dirPath)
		if err = os.MkdirAll(dirPath, 0755); err != nil {
			return errors.Wrapf(err, "mkdirall %s failed", dirPath)
		}

		// Create dst
		var f *os.File
		astilog.Debugf("Creating %s", dst)
		if f, err = os.Create(dst); err != nil {
			err = errors.Wrapf(err, "creating %s failed", dst)
			return
		}
		defer f.Close()

		// Disembed
		var b []byte
		astilog.Debugf("Disembedding %s", src)
		if b, err = p.d(src); err != nil {
			err = errors.Wrapf(err, "disembedding %s failed", src)
			return
		}

		// Copy
		astilog.Debugf("Copying disembedded data to %s", dst)
		if _, err = astiio.Copy(ctx, bytes.NewReader(b), f); err != nil {
			err = errors.Wrapf(err, "copying disembedded data into %s failed", dst)
			return
		}
	} else if err != nil {
		return errors.Wrapf(err, "stating %s failed", dst)
	} else {
		astilog.Debugf("%s already exists, skipping %s disembed...", dst, name)
	}
	return
}
