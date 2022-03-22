package main_test

import (
	"fmt"
	"github.com/asticode/go-astilectron"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestZipShouldRemove(t *testing.T) {
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	appName := "Test"
	dataDir, _ := filepath.Abs("./temp/Test")
	vendorDir := filepath.Join(dataDir, "vendor") // https://github.com/CarsonSlovoka/go-astilectron/blob/e7796e5/paths.go#L56
	astilectronDownloadDst := filepath.Join(vendorDir, fmt.Sprintf("astilectron-v%s.zip", astilectron.DefaultVersionAstilectron))
	astilectronDir := filepath.Join(vendorDir, fmt.Sprintf("astilectron"))
	electronDownloadDst := filepath.Join(vendorDir, fmt.Sprintf("electron-%s-%s-v%s.zip", runtime.GOOS, runtime.GOARCH, astilectron.DefaultVersionElectron))
	electronDir := filepath.Join(vendorDir, fmt.Sprintf("electron-%s-%s", runtime.GOOS, runtime.GOARCH))

	t.Log("Make sure the test directory doesn't exist.")
	if err := os.RemoveAll(dataDir); err != nil && !os.IsNotExist(err) {
		t.FailNow()
	}

	t.Logf("Create a directory for test only:\n%s\n", dataDir)
	if err := os.MkdirAll(dataDir, os.FileMode(0666)); err != nil {
		t.Fatalf(err.Error())
	}

	{
		t.Log("common process")
		a, err := astilectron.New(l, astilectron.Options{
			AppName:           appName,
			DataDirectoryPath: dataDir,
		})
		if err != nil {
			t.Fatalf("main: creating astilectron failed: %s", err.Error())
		}
		defer a.Close()

		a.HandleSignals()

		if err = a.Start(); err != nil {
			t.Fatalf("main: starting astilectron failed: %s", err.Error())
		}

		// Close the app immediately since we care about the file only to avoid access being denied. (delete)
		a.Close()
		time.Sleep(10 * time.Second)
		t.Log("astilectron Close")
	}

	{
		t.Log("Check UnZip successful")
		if _, err := os.Stat(astilectronDir); os.IsNotExist(err) {
			t.Fatalf(err.Error())
		}

		if _, err := os.Stat(electronDir); os.IsNotExist(err) {
			t.Fatalf(err.Error())
		}

		t.Log("Check Zip doesn't exist")
		if _, err := os.Stat(astilectronDownloadDst); !os.IsNotExist(err) {
			t.Fatalf(err.Error())
		}

		if _, err := os.Stat(electronDownloadDst); !os.IsNotExist(err) {
			t.Fatalf(err.Error())
		}
	}

	t.Log("Remove the test directory after the test was done.")
	if err := os.RemoveAll(dataDir); err != nil {
		t.Fatalf(err.Error())
	}
}
