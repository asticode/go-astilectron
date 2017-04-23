package astilectron

import (
	"fmt"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

// Paths represents the set of paths needed by Astilectron
type Paths struct {
	astilectronApplication string
	astilectronDirectory   string
	astilectronDownloadSrc string
	astilectronDownloadDst string
	astilectronUnzipSrc    string
	baseDirectory          string
	electronDirectory      string
	electronDownloadSrc    string
	electronDownloadDst    string
	electronExecutable     string
	electronUnzipSrc       string
	vendorDirectory        string
}

// newPaths creates new paths
func newPaths(baseDirectoryPath string) (p *Paths, err error) {
	// Init base directory path
	p = &Paths{}
	if err = p.initBaseDirectory(baseDirectoryPath); err != nil {
		err = errors.Wrap(err, "initializing base directory failed")
		return
	}

	// Init other paths
	//!\\ Order matters
	p.vendorDirectory = filepath.Join(p.baseDirectory, "vendor")
	p.initAstilectronDirectory()
	p.astilectronApplication = filepath.Join(p.astilectronDirectory, "main.js")
	p.astilectronDownloadSrc = fmt.Sprintf("https://github.com/asticode/astilectron/archive/v%s.zip", versionAstilectron)
	p.astilectronDownloadDst = filepath.Join(p.vendorDirectory, fmt.Sprintf("astilectron-v%s.zip", versionAstilectron))
	p.astilectronUnzipSrc = filepath.Join(p.astilectronDownloadDst, fmt.Sprintf("astilectron-%s", versionAstilectron))
	p.electronDirectory = filepath.Join(p.vendorDirectory, "electron")
	p.initElectronDownloadSrc()
	p.electronDownloadDst = filepath.Join(p.vendorDirectory, fmt.Sprintf("electron-v%s.zip", versionElectron))
	p.electronUnzipSrc = p.electronDownloadDst
	p.initElectronExecutable()
	return
}

// initBaseDirectory initializes the base directory path
func (p *Paths) initBaseDirectory(baseDirectoryPath string) (err error) {
	// No path specified in the options
	p.baseDirectory = baseDirectoryPath
	if len(p.baseDirectory) == 0 {
		// Retrieve current user
		var u *user.User
		if u, err = user.Current(); err != nil {
			err = errors.Wrap(err, "retrieving current user failed")
			return
		}

		// Home directory is empty
		p.baseDirectory = u.HomeDir
		if len(p.baseDirectory) == 0 {
			err = errors.New("home dir path is empty")
			return
		}
	}

	// We need the absolute path
	if p.baseDirectory, err = filepath.Abs(p.baseDirectory); err != nil {
		err = errors.Wrap(err, "computing absolute path failed")
		return
	}
	return
}

// initAstilectronDirectory initializes the astilectron directory path
func (p *Paths) initAstilectronDirectory() {
	if len(*astilectronDirectoryPath) > 0 {
		p.astilectronDirectory = *astilectronDirectoryPath
	} else {
		p.astilectronDirectory = filepath.Join(p.vendorDirectory, "astilectron")
	}
}

// initElectronDownloadSrc initializes the electron download source path
func (p *Paths) initElectronDownloadSrc() {
	switch runtime.GOOS {
	case "darwin":
		p.electronDownloadSrc = fmt.Sprintf("https://github.com/electron/electron/releases/download/v%s/electron-v%s-darwin-x64.zip", versionElectron, versionElectron)
	case "linux":
		p.electronDownloadSrc = fmt.Sprintf("https://github.com/electron/electron/releases/download/v%s/electron-v%s-linux-x64.zip", versionElectron, versionElectron)
	case "windows":
		p.electronDownloadSrc = fmt.Sprintf("https://github.com/electron/electron/releases/download/v%s/electron-v%s-win32-ia32.zip", versionElectron, versionElectron)
	}
}

// initElectronExecutable initializes the electron executable path
func (p *Paths) initElectronExecutable() {
	switch runtime.GOOS {
	case "darwin":
		p.electronExecutable = filepath.Join(p.electronDirectory, "Electron.app", "Contents", "MacOS", "Electron")
	case "linux":
		p.electronExecutable = filepath.Join(p.electronDirectory, "electron")
	case "windows":
		p.electronExecutable = filepath.Join(p.electronDirectory, "electron.exe")
	}
}

// BaseDirectory returns the base directory path
func (p *Paths) BaseDirectory() string {
	return p.baseDirectory
}

// AstilectronApplication returns the astilectron application path
func (p *Paths) AstilectronApplication() string {
	return p.astilectronApplication
}

// AstilectronDirectory returns the astilectron directory path
func (p *Paths) AstilectronDirectory() string {
	return p.astilectronDirectory
}

// AstilectronDownloadDst returns the astilectron download destination path
func (p *Paths) AstilectronDownloadDst() string {
	return p.astilectronDownloadDst
}

// AstilectronDownloadSrc returns the astilectron download source path
func (p *Paths) AstilectronDownloadSrc() string {
	return p.astilectronDownloadSrc
}

// AstilectronUnzipSrc returns the astilectron unzip source path
func (p *Paths) AstilectronUnzipSrc() string {
	return p.astilectronUnzipSrc
}

// ElectronDirectory returns the electron directory path
func (p *Paths) ElectronDirectory() string {
	return p.electronDirectory
}

// ElectronDownloadDst returns the electron download destination path
func (p *Paths) ElectronDownloadDst() string {
	return p.electronDownloadDst
}

// ElectronDownloadSrc returns the electron download source path
func (p *Paths) ElectronDownloadSrc() string {
	return p.electronDownloadSrc
}

// ElectronExecutable returns the electron eecutable path
func (p *Paths) ElectronExecutable() string {
	return p.electronExecutable
}

// ElectronUnzipSrc returns the electron unzip source path
func (p *Paths) ElectronUnzipSrc() string {
	return p.electronUnzipSrc
}

// BaseDirectory returns the vendor directory path
func (p *Paths) VendorDirectory() string {
	return p.vendorDirectory
}
