package astilectron

import (
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"io"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/context"
	"github.com/asticode/go-astitools/exec"
	"github.com/pkg/errors"
)

// Versions
const (
	VersionAstilectron = "0.8.0"
	VersionElectron    = "1.6.5"
)

// Misc vars
var (
	validOSes = map[string]bool{
		"darwin":  true,
		"linux":   true,
		"windows": true,
	}
)

// App event names
const (
	EventNameAppClose      = "app.close"
	EventNameAppCmdStop    = "app.cmd.stop"
	EventNameAppCrash      = "app.crash"
	EventNameAppEventReady = "app.event.ready"
)

// Astilectron represents an object capable of interacting with Astilectron
// TODO Fix race conditions
type Astilectron struct {
	canceller    *asticontext.Canceller
	channelQuit  chan bool
	dispatcher   *Dispatcher
	displayPool  *displayPool
	identifier   *identifier
	options      Options
	paths        *Paths
	provisioner  Provisioner
	reader       *reader
	stderrWriter *astiexec.StdWriter
	writer       *writer
}

// Options represents Astilectron options
type Options struct {
	AppName            string
	AppIconDarwinPath  string // Darwin systems requires a specific .icns file
	AppIconDefaultPath string
	BaseDirectoryPath  string
}

// New creates a new Astilectron instance
func New(o Options) (a *Astilectron, err error) {
	// Validate the OS
	if !IsValidOS(runtime.GOOS) {
		err = errors.Wrapf(err, "OS %s is invalid")
		return
	}

	// Init
	a = &Astilectron{
		canceller:   asticontext.NewCanceller(),
		channelQuit: make(chan bool),
		dispatcher:  newDispatcher(),
		displayPool: newDisplayPool(),
		identifier:  newIdentifier(),
		options:     o,
		provisioner: DefaultProvisioner,
	}

	// Set paths
	if a.paths, err = newPaths(runtime.GOOS, runtime.GOARCH, o); err != nil {
		err = errors.Wrap(err, "creating new paths failed")
		return
	}

	// Add default listeners
	a.On(EventNameAppCmdStop, func(e Event) (deleteListener bool) {
		a.Stop()
		return
	})
	a.On(EventNameDisplayEventAdded, func(e Event) (deleteListener bool) {
		a.displayPool.update(e.Displays)
		return
	})
	a.On(EventNameDisplayEventMetricsChanged, func(e Event) (deleteListener bool) {
		a.displayPool.update(e.Displays)
		return
	})
	a.On(EventNameDisplayEventRemoved, func(e Event) (deleteListener bool) {
		a.displayPool.update(e.Displays)
		return
	})
	return
}

// IsValidOS validates the OS
func IsValidOS(os string) (ok bool) {
	_, ok = validOSes[os]
	return
}

// SetProvisioner sets the provisioner
func (a *Astilectron) SetProvisioner(p Provisioner) *Astilectron {
	a.provisioner = p
	return a
}

// On implements the Listenable interface
func (a *Astilectron) On(eventName string, l Listener) {
	a.dispatcher.addListener(mainTargetID, eventName, l)
}

// Start starts Astilectron
func (a *Astilectron) Start() (err error) {
	// Log
	astilog.Debug("Starting...")

	// Start the dispatcher
	go a.dispatcher.start()

	// Provision
	if err = a.provision(); err != nil {
		return errors.Wrap(err, "provisioning failed")
	}

	// Execute
	if err = a.execute(); err != nil {
		err = errors.Wrap(err, "executing failed")
		return
	}
	return
}

// provision provisions Astilectron
func (a *Astilectron) provision() error {
	astilog.Debug("Provisioning...")
	var ctx, _ = a.canceller.NewContext()
	return a.provisioner.Provision(ctx, *a.dispatcher, a.options.AppName, runtime.GOOS, runtime.GOARCH, *a.paths)
}

// execute executes Astilectron in Electron
func (a *Astilectron) execute() (err error) {
	// Log
	astilog.Debug("Executing...")

	// Create command
	var ctx, _ = a.canceller.NewContext()
	var cmd = exec.CommandContext(ctx, a.paths.AppExecutable(), a.paths.AstilectronApplication())

	// Log stderr
	a.stderrWriter = astiexec.NewStdWriter(func(i []byte) { astilog.Debugf("Stderr says: %s", i) })
	cmd.Stderr = a.stderrWriter

	// Pipe StdIn
	var stdin io.WriteCloser
	if stdin, err = cmd.StdinPipe(); err != nil {
		err = errors.Wrap(err, "piping stdin failed")
		return
	}
	a.writer = newWriter(stdin)

	// Pipe StdOut
	var stdout io.ReadCloser
	if stdout, err = cmd.StdoutPipe(); err != nil {
		err = errors.Wrap(err, "piping stdout failed")
		return
	}

	// Read
	a.reader = newReader(a.dispatcher, stdout)
	go a.reader.read()

	// Execute command
	if err = a.executeCmd(cmd); err != nil {
		return errors.Wrap(err, "executing cmd failed")
	}
	return
}

// executeCmd executes the command
func (a *Astilectron) executeCmd(cmd *exec.Cmd) (err error) {
	var e = synchronousFunc(a.canceller, a, func() {
		// Start command
		astilog.Debugf("Starting cmd %s", strings.Join(cmd.Args, " "))
		if err = cmd.Start(); err != nil {
			err = errors.Wrapf(err, "starting cmd %s failed", strings.Join(cmd.Args, " "))
			return
		}

		// Watch command
		go a.watchCmd(cmd)
	}, EventNameAppEventReady)

	// Update display pool
	if e.Displays != nil {
		a.displayPool.update(e.Displays)
	}
	return
}

// watchCmd watches the cmd execution
func (a *Astilectron) watchCmd(cmd *exec.Cmd) {
	// Wait
	cmd.Wait()

	// Check the canceller to check whether it was a crash
	if !a.canceller.Cancelled() {
		astilog.Debug("App has crashed")
		a.dispatcher.Dispatch(Event{Name: EventNameAppCrash, TargetID: mainTargetID})
	} else {
		astilog.Debug("App has closed")
		a.dispatcher.Dispatch(Event{Name: EventNameAppClose, TargetID: mainTargetID})
	}
	a.dispatcher.Dispatch(Event{Name: EventNameAppCmdStop, TargetID: mainTargetID})
}

// Close closes Astilectron properly
func (a *Astilectron) Close() {
	astilog.Debug("Closing...")
	a.canceller.Cancel()
	a.dispatcher.close()
	if a.reader != nil {
		a.reader.close()
	}
	if a.stderrWriter != nil {
		a.stderrWriter.Close()
	}
	if a.writer != nil {
		a.writer.close()
	}
}

// HandleSignals handles signals
func (a *Astilectron) HandleSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for sig := range ch {
			astilog.Debugf("Received signal %s", sig)
			a.Stop()
		}
	}()
}

// Stop orders Astilectron to stop
func (a *Astilectron) Stop() {
	astilog.Debug("Stopping...")
	a.canceller.Cancel()
	if a.channelQuit != nil {
		close(a.channelQuit)
		a.channelQuit = nil
	}
}

// Wait is a blocking pattern
func (a *Astilectron) Wait() {
	if a.channelQuit == nil {
		return
	}
	<-a.channelQuit
}

// Paths returns the paths
func (a *Astilectron) Paths() Paths {
	return *a.paths
}

// Displays returns the displays
func (a *Astilectron) Displays() []*Display {
	return a.displayPool.all()
}

// PrimaryDisplay returns the primary display
func (a *Astilectron) PrimaryDisplay() *Display {
	return a.displayPool.primary()
}

// NewMenu creates a new app menu
func (a *Astilectron) NewMenu(i []*MenuItemOptions) *Menu {
	return newMenu(nil, a, i, a.canceller, a.dispatcher, a.identifier, a.writer)
}

// NewWindow creates a new window
func (a *Astilectron) NewWindow(url string, o *WindowOptions) (*Window, error) {
	return newWindow(a.options, url, o, a.canceller, a.dispatcher, a.identifier, a.writer)
}

// NewWindowInDisplay creates a new window in a specific display
// This overrides the center attribute
func (a *Astilectron) NewWindowInDisplay(d *Display, url string, o *WindowOptions) (*Window, error) {
	if o.X != nil {
		*o.X += d.Bounds().X
	} else {
		o.X = PtrInt(d.Bounds().X)
	}
	if o.Y != nil {
		*o.Y += d.Bounds().Y
	} else {
		o.Y = PtrInt(d.Bounds().Y)
	}
	return newWindow(a.options, url, o, a.canceller, a.dispatcher, a.identifier, a.writer)
}

// NewTray creates a new tray
func (a *Astilectron) NewTray(o *TrayOptions) *Tray {
	return newTray(o, a.canceller, a.dispatcher, a.identifier, a.writer)
}
