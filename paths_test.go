package astilectron

import (
	"testing"

	"os"

	"path/filepath"

	"github.com/stretchr/testify/assert"
)

func TestPaths(t *testing.T) {
	const k = "APPDATA"
	ad := os.Getenv(k)
	os.Setenv(k, "")
	ep, err := os.Executable()
	ep = filepath.Dir(ep)
	assert.NoError(t, err)
	p, err := newPaths("linux", "amd64", Options{})
	assert.NoError(t, err)
	assert.Equal(t, ep+"/vendor/electron-linux-amd64/electron", p.AppExecutable())
	assert.Equal(t, "", p.AppIconDarwinSrc())
	assert.Equal(t, ep, p.BaseDirectory())
	assert.Equal(t, ep, p.DataDirectory())
	assert.Equal(t, ep+"/vendor/astilectron/main.js", p.AstilectronApplication())
	assert.Equal(t, ep+"/vendor/astilectron", p.AstilectronDirectory())
	assert.Equal(t, ep+"/vendor/astilectron-v"+VersionAstilectron+".zip", p.AstilectronDownloadDst())
	assert.Equal(t, "https://github.com/asticode/astilectron/archive/v"+VersionAstilectron+".zip", p.AstilectronDownloadSrc())
	assert.Equal(t, ep+"/vendor/astilectron-v"+VersionAstilectron+".zip/astilectron-"+VersionAstilectron, p.AstilectronUnzipSrc())
	assert.Equal(t, ep+"/vendor/electron-linux-amd64", p.ElectronDirectory())
	assert.Equal(t, ep+"/vendor/electron-linux-amd64-v"+VersionElectron+".zip", p.ElectronDownloadDst())
	assert.Equal(t, "https://github.com/electron/electron/releases/download/v"+VersionElectron+"/electron-v"+VersionElectron+"-linux-x64.zip", p.ElectronDownloadSrc())
	assert.Equal(t, ep+"/vendor/electron-linux-amd64-v"+VersionElectron+".zip", p.ElectronUnzipSrc())
	assert.Equal(t, ep+"/vendor/status.json", p.ProvisionStatus())
	assert.Equal(t, ep+"/vendor", p.VendorDirectory())
	p, err = newPaths("linux", "", Options{})
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/electron/electron/releases/download/v"+VersionElectron+"/electron-v"+VersionElectron+"-linux-ia32.zip", p.ElectronDownloadSrc())
	p, err = newPaths("linux", "arm", Options{})
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/electron/electron/releases/download/v"+VersionElectron+"/electron-v"+VersionElectron+"-linux-armv7l.zip", p.ElectronDownloadSrc())
	p, err = newPaths("darwin", "", Options{BaseDirectoryPath: "/path/to/base/directory", AppIconDarwinPath: "/path/to/darwin/icon", AppIconDefaultPath: "icon"})
	assert.NoError(t, err)
	assert.Equal(t, "/path/to/base/directory/vendor/electron-darwin-/Electron.app/Contents/MacOS/Electron", p.AppExecutable())
	assert.Equal(t, "/path/to/darwin/icon", p.AppIconDarwinSrc())
	assert.Equal(t, "/path/to/base/directory/icon", p.AppIconDefaultSrc())
	assert.Equal(t, "https://github.com/electron/electron/releases/download/v"+VersionElectron+"/electron-v"+VersionElectron+"-darwin-x64.zip", p.ElectronDownloadSrc())
	p, err = newPaths("darwin", "amd64", Options{AppName: "Test app", BaseDirectoryPath: "/path/to/base/directory", DataDirectoryPath: "/path/to/data/directory"})
	assert.NoError(t, err)
	assert.Equal(t, "/path/to/data/directory", p.DataDirectory())
	assert.Equal(t, "/path/to/data/directory/vendor/electron-darwin-amd64/Test app.app/Contents/MacOS/Test app", p.AppExecutable())
	assert.Equal(t, "/path/to/data/directory/vendor/electron-darwin-amd64-v"+VersionElectron+".zip", p.ElectronDownloadDst())
	assert.Equal(t, "/path/to/data/directory/vendor/electron-darwin-amd64-v"+VersionElectron+".zip", p.ElectronUnzipSrc())
	const pad = "/path/to/appdata"
	os.Setenv(k, pad)
	p, err = newPaths("windows", "amd64", Options{})
	assert.NoError(t, err)
	assert.Equal(t, pad, p.DataDirectory())
	assert.Equal(t, pad+"/vendor", p.VendorDirectory())
	assert.Equal(t, pad+"/vendor/electron-windows-amd64/electron.exe", p.AppExecutable())
	assert.Equal(t, "https://github.com/electron/electron/releases/download/v"+VersionElectron+"/electron-v"+VersionElectron+"-win32-x64.zip", p.ElectronDownloadSrc())
	assert.Equal(t, pad+"/vendor/electron-windows-amd64-v"+VersionElectron+".zip", p.ElectronDownloadDst())
	assert.Equal(t, pad+"/vendor/electron-windows-amd64-v"+VersionElectron+".zip", p.ElectronUnzipSrc())
	p, err = newPaths("windows", "", Options{})
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/electron/electron/releases/download/v"+VersionElectron+"/electron-v"+VersionElectron+"-win32-ia32.zip", p.ElectronDownloadSrc())
	os.Setenv(k, ad)
}
