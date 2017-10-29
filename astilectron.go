package astilectron

import (
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/context"
	"github.com/asticode/go-astitools/exec"
	"github.com/pkg/errors"
)

// Versions
const (
	VersionAstilectron = "0.10.0"
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
	EventNameAppClose         = "app.close"
	EventNameAppCmdStop       = "app.cmd.stop"
	EventNameAppCrash         = "app.crash"
	EventNameAppErrorAccept   = "app.error.accept"
	EventNameAppEventReady    = "app.event.ready"
	EventNameAppNoAccept      = "app.no.accept"
	EventNameAppTooManyAccept = "app.too.many.accept"
)

// Astilectron represents an object capable of interacting with Astilectron
// TODO Fix race conditions
type Astilectron struct {
	canceller    *asticontext.Canceller
	channelQuit  chan bool
	closeOnce    sync.Once
	dispatcher   *Dispatcher
	displayPool  *displayPool
	identifier   *identifier
	listener     net.Listener
	options      Options
	paths        *Paths
	provisioner  Provisioner
	reader       *reader
	stderrWriter *astiexec.StdWriter
	stdoutWriter *astiexec.StdWriter
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

	// Unfortunately communicating with Electron through stdin/stdout doesn't work on Windows so all communications
	// will be done through TCP
	if err = a.listenTCP(); err != nil {
		return errors.Wrap(err, "listening failed")
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
	return a.provisioner.Provision(ctx, a.dispatcher, a.options.AppName, runtime.GOOS, runtime.GOARCH, *a.paths)
}

// listenTCP listens to the first TCP connection coming its way (this should be Astilectron)
func (a *Astilectron) listenTCP() (err error) {
	// Log
	astilog.Debug("Listening...")

	// Listen
	if a.listener, err = net.Listen("tcp", "127.0.0.1:"); err != nil {
		return errors.Wrap(err, "tcp net.Listen failed")
	}

	// Check a connection has been accepted quickly enough
	var chanAccepted = make(chan bool)
	go a.watchNoAccept(30*time.Second, chanAccepted)

	// Accept connections
	go a.acceptTCP(chanAccepted)
	return
}

// watchNoAccept checks whether a TCP connection is accepted quickly enough
func (a *Astilectron) watchNoAccept(timeout time.Duration, chanAccepted chan bool) {
	var t = time.NewTimer(timeout)
	defer t.Stop()
	for {
		select {
		case <-chanAccepted:
			return
		case <-t.C:
			astilog.Errorf("No TCP connection has been accepted in the past %s", timeout)
			a.dispatcher.Dispatch(Event{Name: EventNameAppNoAccept, TargetID: mainTargetID})
			a.dispatcher.Dispatch(Event{Name: EventNameAppCmdStop, TargetID: mainTargetID})
			return
		}
	}
}

// watchAcceptTCP accepts TCP connections
func (a *Astilectron) acceptTCP(chanAccepted chan bool) {
	for i := 0; i <= 1; i++ {
		// Accept
		var conn net.Conn
		var err error
		if conn, err = a.listener.Accept(); err != nil {
			astilog.Errorf("%s while TCP accepting", err)
			a.dispatcher.Dispatch(Event{Name: EventNameAppErrorAccept, TargetID: mainTargetID})
			a.dispatcher.Dispatch(Event{Name: EventNameAppCmdStop, TargetID: mainTargetID})
			return
		}

		// We only accept the first connection which should be Astilectron, close the next one and stop
		// the app
		if i > 0 {
			astilog.Errorf("Too many TCP connections")
			a.dispatcher.Dispatch(Event{Name: EventNameAppTooManyAccept, TargetID: mainTargetID})
			a.dispatcher.Dispatch(Event{Name: EventNameAppCmdStop, TargetID: mainTargetID})
			conn.Close()
			return
		}

		// Let the timer know a connection has been accepted
		chanAccepted <- true

		// Create reader and writer
		a.writer = newWriter(conn)
		a.reader = newReader(a.dispatcher, conn)
		go a.reader.read()
	}
}

// execute executes Astilectron in Electron
func (a *Astilectron) execute() (err error) {
	// Log
	astilog.Debug("Executing...")

	// Create command
	var ctx, _ = a.canceller.NewContext()
	var cmd = exec.CommandContext(ctx, a.paths.AppExecutable(), a.paths.AstilectronApplication(), a.listener.Addr().String())
	a.stderrWriter = astiexec.NewStdWriter(func(i []byte) { astilog.Debugf("Stderr says: %s", i) })
	a.stdoutWriter = astiexec.NewStdWriter(func(i []byte) { astilog.Debugf("Stdout says: %s", i) })
	cmd.Stderr = a.stderrWriter
	cmd.Stdout = a.stdoutWriter

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
	if a.listener != nil {
		a.listener.Close()
	}
	if a.reader != nil {
		a.reader.close()
	}
	if a.stderrWriter != nil {
		a.stderrWriter.Close()
	}
	if a.stdoutWriter != nil {
		a.stdoutWriter.Close()
	}
	if a.writer != nil {
		a.writer.close()
	}
}

// HandleSignals handles signals
func (a *Astilectron) HandleSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGABRT, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
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
	a.closeOnce.Do(func() {
		close(a.channelQuit)
	})
}

// Wait is a blocking pattern
func (a *Astilectron) Wait() {
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
	return newMenu(nil, mainTargetID, i, a.canceller, a.dispatcher, a.identifier, a.writer)
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
