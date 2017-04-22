package astilectron

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/context"
	"github.com/asticode/go-astitools/slice"
	"github.com/pkg/errors"
)

// Constants
const (
	defaultApplicationName = "Astilectron"
	versionAstilectron     = "0.1.0"
	versionElectron        = "1.6.5"
)

// Vars
var (
	astilectronDirectoryPath = flag.String("astilectron-directory-path", "", "the astilectron directory path")
	validOSes                = []string{"darwin", "linux", "windows"}
)

// Astilectron represents an object capable of interacting with Astilectron
type Astilectron struct {
	applicationName string
	canceller       *asticontext.Canceller
	channelQuit     chan bool
	paths           *Paths
	provisioner     Provisioner
}

// Options represents Astilectron options
type Options struct {
	ApplicationName   string
	BaseDirectoryPath string
}

// New creates a new Astilectron instance
func New(o Options) (a *Astilectron, err error) {
	// Validate the OS
	if err = validateOS(); err != nil {
		err = errors.Wrap(err, "validating OS failed")
		return
	}
	a = &Astilectron{
		canceller:   asticontext.NewCanceller(),
		channelQuit: make(chan bool),
		provisioner: DefaultProvisioner,
	}

	// Set application name
	a.applicationName = defaultApplicationName
	if len(o.ApplicationName) > 0 {
		a.applicationName = o.ApplicationName
	}

	// Set paths
	if a.paths, err = newPaths(o.BaseDirectoryPath); err != nil {
		err = errors.Wrap(err, "creating new paths failed")
		return
	}
	return
}

// validateOS validates the OS
func validateOS() error {
	if !astislice.InStringSlice(runtime.GOOS, validOSes) {
		return fmt.Errorf("OS %s is not supported", runtime.GOOS)
	}
	return nil
}

// SetProvisioner sets the provisioner
func (a *Astilectron) SetProvisioner(p Provisioner) *Astilectron {
	a.provisioner = p
	return a
}

// Start starts Astilectron
func (a *Astilectron) Start() (err error) {
	// Log
	astilog.Debug("Starting...")

	// TODO Start handling events so that we can send provision.start, download.start, unzip.start, provision.done,
	// etc. events

	// Provision
	if err = a.provision(); err != nil {
		return errors.Wrap(err, "provisioning failed")
	}

	// Execute
	var stdin io.WriteCloser
	var stdout io.ReadCloser
	if stdin, stdout, err = a.execute(); err != nil {
		err = errors.Wrap(err, "executing failed")
		return
	}
	_, _ = stdin, stdout
	return
}

// provision provisions Astilectron
func (a *Astilectron) provision() error {
	astilog.Debug("Provisioning...")
	return a.provisioner.Provision(a.paths)
}

// execute executes Astilectron in Electron
func (a *Astilectron) execute() (stdin io.WriteCloser, stdout io.ReadCloser, err error) {
	// Log
	astilog.Debug("Executing...")

	// Create command
	var ctx, _ = a.canceller.NewContext()
	var cmd = exec.CommandContext(ctx, a.paths.ElectronExecutable(), a.paths.AstilectronApplication())

	// Pipe StdIn
	if stdin, err = cmd.StdinPipe(); err != nil {
		err = errors.Wrap(err, "piping stdin failed")
		return
	}

	// Pipe StdOut
	if stdout, err = cmd.StdoutPipe(); err != nil {
		err = errors.Wrap(err, "piping stdout failed")
		return
	}

	// Start command
	astilog.Debugf("Starting cmd %s", strings.Join(cmd.Args, " "))
	if err = cmd.Start(); err != nil {
		err = errors.Wrapf(err, "starting cmd %s failed", strings.Join(cmd.Args, " "))
		return
	}
	return
}

// Close closes Astilectron properly
func (t *Astilectron) Close() {
	astilog.Debug("Closing...")
	t.canceller.Cancel()
}

// HandleSignals handles signals
func (t *Astilectron) HandleSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for sig := range ch {
			astilog.Debugf("Received signal %s", sig)
			t.stop()
		}
	}()
}

// Stop orders Astilectron to stop
func (t *Astilectron) stop() {
	astilog.Debug("Stopping...")
	close(t.channelQuit)
}

// Wait is a blocking pattern
func (t *Astilectron) Wait() {
	for {
		select {
		case <-t.channelQuit:
			return
		}
	}
}
