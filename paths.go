package astilectron

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Paths represents the set of paths needed by Astilectron
type Paths struct {
	appExecutable          string
	appIconDarwinSrc       string
	astilectronApplication string
	astilectronDirectory   string
	astilectronDownloadSrc string
	astilectronDownloadDst string
	astilectronUnzipSrc    string
	baseDirectory          string
	electronDirectory      string
	electronDownloadSrc    string
	electronDownloadDst    string
	electronUnzipSrc       string
	provisionStatus        string
	vendorDirectory        string
}

// newPaths creates new paths
func newPaths(os, arch string, o Options) (p *Paths, err error) {
	// Init base directory path
	p = &Paths{}
	if err = p.initBaseDirectory(o.BaseDirectoryPath); err != nil {
		err = errors.Wrap(err, "initializing base directory failed")
		return
	}

	// Init other paths
	//!\\ Order matters
	p.appIconDarwinSrc = o.AppIconDarwinPath
	p.vendorDirectory = filepath.Join(p.baseDirectory, "vendor")
	p.provisionStatus = filepath.Join(p.vendorDirectory, "status.json")
	p.initAstilectronDirectory()
	p.astilectronApplication = filepath.Join(p.astilectronDirectory, "main.js")
	p.astilectronDownloadSrc = AstilectronDownloadSrc()
	p.astilectronDownloadDst = filepath.Join(p.vendorDirectory, fmt.Sprintf("astilectron-v%s.zip", VersionAstilectron))
	p.astilectronUnzipSrc = filepath.Join(p.astilectronDownloadDst, fmt.Sprintf("astilectron-%s", VersionAstilectron))
	p.electronDirectory = filepath.Join(p.vendorDirectory, "electron")
	p.electronDownloadSrc = ElectronDownloadSrc(os, arch)
	// TODO Split folders for each OS/ARCH couples?
	p.electronDownloadDst = filepath.Join(p.vendorDirectory, fmt.Sprintf("electron-v%s.zip", VersionElectron))
	p.electronUnzipSrc = p.electronDownloadDst
	p.initAppExecutable(os, o.AppName)
	return
}

// initBaseDirectory initializes the base directory path
func (p *Paths) initBaseDirectory(baseDirectoryPath string) (err error) {
	// No path specified in the options
	p.baseDirectory = baseDirectoryPath
	if len(p.baseDirectory) == 0 {
		// Retrieve executable path
		var ep string
		if ep, err = os.Executable(); err != nil {
			err = errors.Wrap(err, "retrieving executable path failed")
			return
		}
		p.baseDirectory = filepath.Dir(ep)
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

// AstilectronDownloadSrc returns the download URL of the (currently platform-independant) astilectron zipfile
func AstilectronDownloadSrc() string {
	return fmt.Sprintf("https://github.com/asticode/astilectron/archive/v%s.zip", VersionAstilectron)
}

// ElectronDownloadSrc returns the download URL of the platform-dependant electron zipfile
func ElectronDownloadSrc(os, arch string) string {
	// Get OS name
	var o string
	switch strings.ToLower(os) {
	case "darwin":
		o = "darwin"
	case "linux":
		o = "linux"
	case "windows":
		o = "win32"
	}

	// Get arch name
	var a = "ia32"
	if strings.ToLower(arch) == "amd64" || o == "darwin" {
		a = "x64"
	}

	// Return url
	return fmt.Sprintf("https://github.com/electron/electron/releases/download/v%s/electron-v%s-%s-%s.zip", VersionElectron, VersionElectron, o, a)
}

// initAppExecutable initializes the app executable path
func (p *Paths) initAppExecutable(os, appName string) {
	switch os {
	case "darwin":
		if appName == "" {
			appName = "Electron"
		}
		p.appExecutable = filepath.Join(p.electronDirectory, appName+".app", "Contents", "MacOS", appName)
	case "linux":
		p.appExecutable = filepath.Join(p.electronDirectory, "electron")
	case "windows":
		p.appExecutable = filepath.Join(p.electronDirectory, "electron.exe")
	}
}

// AppExecutable returns the app executable path
func (p *Paths) AppExecutable() string {
	return p.appExecutable
}

// AppIconDarwinSrc returns the darwin app icon path
func (p *Paths) AppIconDarwinSrc() string {
	return p.appIconDarwinSrc
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

// ElectronUnzipSrc returns the electron unzip source path
func (p *Paths) ElectronUnzipSrc() string {
	return p.electronUnzipSrc
}

// ProvisionStatus returns the provision status path
func (p *Paths) ProvisionStatus() string {
	return p.provisionStatus
}

// VendorDirectory returns the vendor directory path
func (p *Paths) VendorDirectory() string {
	return p.vendorDirectory
}
