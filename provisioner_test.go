package astilectron

import (
	"context"
	"github.com/asticode/go-astikit"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func testProvisionerSuccessful(t *testing.T, p Paths, osName, arch, versionAstilectron, versionElectron string) {
	_, err := os.Stat(p.AstilectronApplication())
	assert.NoError(t, err)
	_, err = os.Stat(p.AppExecutable())
	assert.NoError(t, err)
	b, err := ioutil.ReadFile(p.ProvisionStatus())
	assert.NoError(t, err)
	assert.Equal(t, "{\"astilectron\":{\"version\":\""+versionAstilectron+"\"},\"electron\":{\""+provisionStatusElectronKey(osName, arch)+"\":{\"version\":\""+versionElectron+"\"}}}\n", string(b))
}

func TestDefaultProvisioner(t *testing.T) {
	// Init
	var o = Options{BaseDirectoryPath: mockedTempPath()}
	defer os.RemoveAll(o.BaseDirectoryPath)
	var mh = &mockedHandler{}
	var s = httptest.NewServer(mh)

	// Test linux
	p, err := newPaths("linux", "amd64", o)
	assert.NoError(t, err)
	p.astilectronUnzipSrc = filepath.Join(p.astilectronDownloadDst, "astilectron")
	p.astilectronDownloadSrc = s.URL + "/provisioner/astilectron"
	p.electronDownloadSrc = s.URL + "/provisioner/electron/linux"
	err = newDefaultProvisioner(nil).Provision(context.Background(), "", "linux", "amd64", DefaultVersionAstilectron, DefaultVersionElectron, *p)
	assert.NoError(t, err)
	testProvisionerSuccessful(t, *p, "linux", "amd64", DefaultVersionAstilectron, DefaultVersionElectron)

	// Test nothing happens if provision status is up to date
	mh.e = true
	os.Remove(p.AstilectronDownloadDst())
	os.Remove(p.ElectronDownloadDst())
	err = newDefaultProvisioner(nil).Provision(context.Background(), "", "linux", "amd64", DefaultVersionAstilectron, DefaultVersionElectron, *p)
	assert.NoError(t, err)
	testProvisionerSuccessful(t, *p, "linux", "amd64", DefaultVersionAstilectron, DefaultVersionElectron)

	// Test windows
	mh.e = false
	os.RemoveAll(o.BaseDirectoryPath)
	p, err = newPaths("windows", "amd64", o)
	assert.NoError(t, err)
	p.astilectronUnzipSrc = filepath.Join(p.astilectronDownloadDst, "astilectron")
	p.astilectronDownloadSrc = s.URL + "/provisioner/astilectron"
	p.electronDownloadSrc = s.URL + "/provisioner/electron/windows"
	err = newDefaultProvisioner(nil).Provision(context.Background(), "", "windows", "amd64", DefaultVersionAstilectron, DefaultVersionElectron, *p)
	assert.NoError(t, err)
	testProvisionerSuccessful(t, *p, "windows", "amd64", DefaultVersionAstilectron, DefaultVersionElectron)

	// Test darwin without custom app name + icon
	os.RemoveAll(o.BaseDirectoryPath)
	p, err = newPaths("darwin", "amd64", o)
	assert.NoError(t, err)
	p.astilectronUnzipSrc = filepath.Join(p.astilectronDownloadDst, "astilectron")
	p.astilectronDownloadSrc = s.URL + "/provisioner/astilectron"
	p.electronDownloadSrc = s.URL + "/provisioner/electron/darwin"
	err = newDefaultProvisioner(nil).Provision(context.Background(), "", "darwin", "amd64", DefaultVersionAstilectron, DefaultVersionElectron, *p)
	assert.NoError(t, err)
	testProvisionerSuccessful(t, *p, "darwin", "amd64", DefaultVersionAstilectron, DefaultVersionElectron)

	// Test darwin with custom app name + icon
	os.RemoveAll(o.BaseDirectoryPath)
	o.AppName = "Test app"
	wd, err := os.Getwd()
	assert.NoError(t, err)
	o.AppIconDarwinPath = filepath.Join(wd, "testdata", "provisioner", "icon.icns")
	p, err = newPaths("darwin", "amd64", o)
	assert.NoError(t, err)
	p.astilectronUnzipSrc = filepath.Join(p.astilectronDownloadDst, "astilectron")
	p.astilectronDownloadSrc = s.URL + "/provisioner/astilectron"
	p.electronDownloadSrc = s.URL + "/provisioner/electron/darwin"
	err = newDefaultProvisioner(nil).Provision(context.Background(), o.AppName, "darwin", "amd64", DefaultVersionAstilectron, DefaultVersionElectron, *p)
	assert.NoError(t, err)
	testProvisionerSuccessful(t, *p, "darwin", "amd64", DefaultVersionAstilectron, DefaultVersionElectron)
	// Rename
	_, err = os.Stat(filepath.Join(p.ElectronDirectory(), o.AppName+".app"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(p.ElectronDirectory(), o.AppName+".app", "Contents", "MacOS", o.AppName))
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
	b, err = ioutil.ReadFile(filepath.Join(p.ElectronDirectory(), o.AppName+".app", "Contents", "Frameworks", o.AppName+" Helper.app", "Contents", "Info.plist"))
	assert.NoError(t, err)
	assert.Equal(t, "<string>"+o.AppName+" Test</string>", string(b))
}

func TestNewDisembedderProvisioner(t *testing.T) {
	// Init
	var o = Options{BaseDirectoryPath: mockedTempPath()}
	defer os.RemoveAll(o.BaseDirectoryPath)
	p, err := newPaths("linux", "amd64", o)
	assert.NoError(t, err)
	p.astilectronUnzipSrc = filepath.Join(p.astilectronDownloadDst, "astilectron-0.35.1")
	pvb := NewDisembedderProvisioner(mockedDisembedder, "astilectron", "electron/linux", nil)

	// Test provision
	err = pvb.Provision(context.Background(), "", "linux", "amd64", DefaultVersionAstilectron, DefaultVersionElectron, *p)
	assert.NoError(t, err)
	testProvisionerSuccessful(t, *p, "linux", "amd64", DefaultVersionAstilectron, DefaultVersionElectron)
}

func TestRemoveDownloadDst(t *testing.T) {
	var o = Options{
		DataDirectoryPath: mockedTempPath(),
	}

	// Make sure the test directory doesn't exist.
	if err := os.RemoveAll(o.DataDirectoryPath); err != nil && !os.IsNotExist(err) {
		t.Fatalf("main: removing %s failed: %s", o.DataDirectoryPath, err)
	}
	defer os.RemoveAll(o.DataDirectoryPath)

	a, err := New(astikit.AdaptTestLogger(t), o)
	if err != nil {
		t.Fatalf("main: creating astilectron failed: %s", err)
	}

	p := a.Paths()

	if err = a.provision(); err != nil {
		t.Fatalf("main: provisionning failed: %s", err)
	}

	// Check UnZip successful
	if _, err := os.Stat(p.AstilectronDirectory()); os.IsNotExist(err) {
		t.Fatalf("%v", err)
	}

	if _, err := os.Stat(p.ElectronDirectory()); os.IsNotExist(err) {
		t.Fatalf("%v", err)
	}

	// Check Zip doesn't exist
	if _, err := os.Stat(p.AstilectronDownloadDst()); !os.IsNotExist(err) {
		t.Fatalf("%v", err)
	}

	if _, err := os.Stat(p.ElectronDownloadDst()); !os.IsNotExist(err) {
		t.Fatalf("%v", err)
	}
}
