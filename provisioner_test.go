package astilectron

import (
	"context"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testProvisionerSuccessful(t *testing.T, p Paths) {
	_, err := os.Stat(p.AstilectronApplication())
	assert.NoError(t, err)
	_, err = os.Stat(p.AppExecutable())
	assert.NoError(t, err)
	b, err := ioutil.ReadFile(p.ProvisionStatus())
	assert.NoError(t, err)
	assert.Equal(t, "{\"astilectron\":{\"version\":\""+versionAstilectron+"\"},\"electron\":{\"version\":\""+versionElectron+"\"}}\n", string(b))
}

func TestDefaultProvisioner(t *testing.T) {
	// Init
	var o = Options{BaseDirectoryPath: mockedTempPath()}
	defer os.RemoveAll(o.BaseDirectoryPath)
	var mh = &mockedHandler{}
	var s = httptest.NewServer(mh)
	var d = newDispatcher()
	defer d.close()
	go d.start()

	// Test linux
	p, err := newPaths("linux", o)
	assert.NoError(t, err)
	p.astilectronUnzipSrc = filepath.Join(p.astilectronDownloadDst, "astilectron")
	p.astilectronDownloadSrc = s.URL + "/provisioner/astilectron"
	p.electronDownloadSrc = s.URL + "/provisioner/electron/linux"
	err = DefaultProvisioner.Provision(context.Background(), *d, "", "linux", *p)
	assert.NoError(t, err)
	testProvisionerSuccessful(t, *p)

	// Test nothing happens if provision status is up to date
	mh.e = true
	os.Remove(p.AstilectronDownloadDst())
	os.Remove(p.ElectronDownloadDst())
	err = DefaultProvisioner.Provision(context.Background(), *d, "", "linux", *p)
	assert.NoError(t, err)
	testProvisionerSuccessful(t, *p)

	// Test windows
	mh.e = false
	os.RemoveAll(o.BaseDirectoryPath)
	p, err = newPaths("windows", o)
	assert.NoError(t, err)
	p.astilectronUnzipSrc = filepath.Join(p.astilectronDownloadDst, "astilectron")
	p.astilectronDownloadSrc = s.URL + "/provisioner/astilectron"
	p.electronDownloadSrc = s.URL + "/provisioner/electron/windows"
	err = DefaultProvisioner.Provision(context.Background(), *d, "", "windows", *p)
	assert.NoError(t, err)
	testProvisionerSuccessful(t, *p)

	// Test darwin without custom app name + icon
	os.RemoveAll(o.BaseDirectoryPath)
	p, err = newPaths("darwin", o)
	assert.NoError(t, err)
	p.astilectronUnzipSrc = filepath.Join(p.astilectronDownloadDst, "astilectron")
	p.astilectronDownloadSrc = s.URL + "/provisioner/astilectron"
	p.electronDownloadSrc = s.URL + "/provisioner/electron/darwin"
	err = DefaultProvisioner.Provision(context.Background(), *d, "", "darwin", *p)
	assert.NoError(t, err)
	testProvisionerSuccessful(t, *p)

	// Test darwin with custom app name + icon
	os.RemoveAll(o.BaseDirectoryPath)
	o.AppName = "Test app"
	o.AppIconDarwinPath = "testdata/provisioner/icon.icns"
	p, err = newPaths("darwin", o)
	assert.NoError(t, err)
	p.astilectronUnzipSrc = filepath.Join(p.astilectronDownloadDst, "astilectron")
	p.astilectronDownloadSrc = s.URL + "/provisioner/astilectron"
	p.electronDownloadSrc = s.URL + "/provisioner/electron/darwin"
	err = DefaultProvisioner.Provision(context.Background(), *d, o.AppName, "darwin", *p)
	assert.NoError(t, err)
	testProvisionerSuccessful(t, *p)
	// Rename
	_, err = os.Stat(filepath.Join(p.ElectronDirectory(), o.AppName+".app"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(p.ElectronDirectory(), o.AppName+".app", "Contents", "MacOS", o.AppName))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(p.ElectronDirectory(), o.AppName+".app", "Contents", "Frameworks", o.AppName+" Helper EH.app"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(p.ElectronDirectory(), o.AppName+".app", "Contents", "Frameworks", o.AppName+" Helper EH.app", "Contents", "MacOS", o.AppName+" Helper EH"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(p.ElectronDirectory(), o.AppName+".app", "Contents", "Frameworks", o.AppName+" Helper NP.app"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(p.ElectronDirectory(), o.AppName+".app", "Contents", "Frameworks", o.AppName+" Helper NP.app", "Contents", "MacOS", o.AppName+" Helper NP"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(p.ElectronDirectory(), o.AppName+".app", "Contents", "Frameworks", o.AppName+" Helper.app"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(p.ElectronDirectory(), o.AppName+".app", "Contents", "Frameworks", o.AppName+" Helper.app", "Contents", "MacOS", o.AppName+" Helper"))
	assert.NoError(t, err)
	// Icon
	b, err := ioutil.ReadFile(filepath.Join(p.ElectronDirectory(), o.AppName+".app", "Contents", "Resources", "electron.icns"))
	assert.NoError(t, err)
	assert.Equal(t, "body", string(b))
	// Replace
	b, err = ioutil.ReadFile(filepath.Join(p.ElectronDirectory(), o.AppName+".app", "Contents", "Info.plist"))
	assert.NoError(t, err)
	assert.Equal(t, "<string>"+o.AppName+" Test</string>", string(b))
	b, err = ioutil.ReadFile(filepath.Join(p.ElectronDirectory(), o.AppName+".app", "Contents", "Frameworks", o.AppName+" Helper EH.app", "Contents", "Info.plist"))
	assert.NoError(t, err)
	assert.Equal(t, "<string>"+o.AppName+" Test</string>", string(b))
	b, err = ioutil.ReadFile(filepath.Join(p.ElectronDirectory(), o.AppName+".app", "Contents", "Frameworks", o.AppName+" Helper NP.app", "Contents", "Info.plist"))
	assert.NoError(t, err)
	assert.Equal(t, "<string>"+o.AppName+" Test</string>", string(b))
	b, err = ioutil.ReadFile(filepath.Join(p.ElectronDirectory(), o.AppName+".app", "Contents", "Frameworks", o.AppName+" Helper.app", "Contents", "Info.plist"))
	assert.NoError(t, err)
	assert.Equal(t, "<string>"+o.AppName+" Test</string>", string(b))
}

func TestNewDisembedderProvisioner(t *testing.T) {
	// Init
	var o = Options{BaseDirectoryPath: mockedTempPath()}
	defer os.RemoveAll(o.BaseDirectoryPath)
	var d = newDispatcher()
	defer d.close()
	go d.start()
	p, err := newPaths("linux", o)
	assert.NoError(t, err)
	p.astilectronUnzipSrc = filepath.Join(p.astilectronDownloadDst, "astilectron")
	pvb := NewDisembedderProvisioner(mockedDisembedder, "astilectron", "electron/linux")

	// Test provision
	err = pvb.Provision(context.Background(), *d, "", "linux", *p)
	assert.NoError(t, err)
	testProvisionerSuccessful(t, *p)
}
