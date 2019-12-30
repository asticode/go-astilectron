package astilectron

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"time"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Versions
const (
	DefaultAcceptTCPTimeout   = 30 * time.Second
	DefaultVersionAstilectron = "0.34.0"
	DefaultVersionElectron    = "4.0.1"
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
	EventNameAppCmdQuit       = "app.cmd.quit" // Sends an event to Electron to properly quit the app
	EventNameAppCmdStop       = "app.cmd.stop" // Cancel the context which results in exiting abruptly Electron's app
	EventNameAppCrash         = "app.crash"
	EventNameAppErrorAccept   = "app.error.accept"
	EventNameAppEventReady    = "app.event.ready"
	EventNameAppNoAccept      = "app.no.accept"
	EventNameAppTooManyAccept = "app.too.many.accept"
)

// Astilectron represents an object capable of interacting with Astilectron
type Astilectron struct {
	dispatcher   *dispatcher
	displayPool  *displayPool
	dock         *Dock
	executer     Executer
	identifier   *identifier
	listener     net.Listener
	options      Options
	paths        *Paths
	provisioner  Provisioner
	reader       *reader
	stderrWriter *astikit.WriterAdapter
	stdoutWriter *astikit.WriterAdapter
	supported    *Supported
	worker       *astikit.Worker
	writer       *writer
}

// Options represents Astilectron options
type Options struct {
	AcceptTCPTimeout   time.Duration
	AppName            string
	AppIconDarwinPath  string // Darwin systems requires a specific .icns file
	AppIconDefaultPath string
	BaseDirectoryPath  string
	DataDirectoryPath  string
	ElectronSwitches   []string
	SingleInstance     bool
	SkipSetup          bool // If true, the user must handle provisioning and executing astilectron.
	TCPPort            *int // The port to listen on.
	VersionAstilectron string
	VersionElectron    string
}

// Supported represents Astilectron supported features
type Supported struct {
	Notification *bool `json:"notification"`
}

// New creates a new Astilectron instance
func New(o Options) (a *Astilectron, err error) {
	// Validate the OS
	if !IsValidOS(runtime.GOOS) {
		err = errors.Wrapf(err, "OS %s is invalid", runtime.GOOS)
		return
	}

	if o.VersionAstilectron == "" {
		o.VersionAstilectron = DefaultVersionAstilectron
	}
	if o.VersionElectron == "" {
		o.VersionElectron = DefaultVersionElectron
	}

	// Init
	a = &Astilectron{
		dispatcher:  newDispatcher(),
		displayPool: newDisplayPool(),
		executer:    DefaultExecuter,
		identifier:  newIdentifier(),
		options:     o,
		provisioner: newDefaultProvisioner(astilog.GetLogger()),
		worker:      astikit.NewWorker(astikit.WorkerOptions{Logger: astilog.GetLogger()}),
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
	a.On(EventNameAppCmdQuit, func(e Event) (deleteListener bool) {
		a.Stop()
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

// SetExecuter sets the executer
func (a *Astilectron) SetExecuter(e Executer) *Astilectron {
	a.executer = e
	return a
}

// On implements the Listenable interface
func (a *Astilectron) On(eventName string, l Listener) {
	a.dispatcher.addListener(targetIDApp, eventName, l)
}

// Start starts Astilectron
func (a *Astilectron) Start() (err error) {
	// Log
	astilog.Debug("Starting...")

	// Provision
	if !a.options.SkipSetup {
		if err = a.provision(); err != nil {
			return errors.Wrap(err, "provisioning failed")
		}
	}

	// Unfortunately communicating with Electron through stdin/stdout doesn't work on Windows so all communications
	// will be done through TCP
	if err = a.listenTCP(); err != nil {
		return errors.Wrap(err, "listening failed")
	}

	// Execute
	if !a.options.SkipSetup {
		if err = a.execute(); err != nil {
			return errors.Wrap(err, "executing failed")
		}
	} else {
		synchronousFunc(a.worker.Context(), a, nil, "app.event.ready")
	}
	return nil
}

// provision provisions Astilectron
func (a *Astilectron) provision() error {
	astilog.Debug("Provisioning...")
	return a.provisioner.Provision(a.worker.Context(), a.options.AppName, runtime.GOOS, runtime.GOARCH, a.options.VersionAstilectron, a.options.VersionElectron, *a.paths)
}

// listenTCP creates a TCP server for astilectron to connect to
// and listens to the first TCP connection coming its way (this should be Astilectron).
func (a *Astilectron) listenTCP() (err error) {
	// Log
	astilog.Debug("Listening...")

	addr := "127.0.0.1:"
	if a.options.TCPPort != nil {
		addr += fmt.Sprint(*a.options.TCPPort)
	}
	// Listen
	if a.listener, err = net.Listen("tcp", addr); err != nil {
		return errors.Wrap(err, "tcp net.Listen failed")
	}

	// Check a connection has been accepted quickly enough
	var chanAccepted = make(chan bool)
	go a.watchNoAccept(a.options.AcceptTCPTimeout, chanAccepted)

	// Accept connections
	go a.acceptTCP(chanAccepted)
	return
}

// watchNoAccept checks whether a TCP connection is accepted quickly enough
func (a *Astilectron) watchNoAccept(timeout time.Duration, chanAccepted chan bool) {
	//check timeout
	if timeout == 0 {
		timeout = DefaultAcceptTCPTimeout
	}
	var t = time.NewTimer(timeout)
	defer t.Stop()
	for {
		select {
		case <-chanAccepted:
			return
		case <-t.C:
			astilog.Errorf("No TCP connection has been accepted in the past %s", timeout)
			a.dispatcher.dispatch(Event{Name: EventNameAppNoAccept, TargetID: targetIDApp})
			a.dispatcher.dispatch(Event{Name: EventNameAppCmdStop, TargetID: targetIDApp})
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
			a.dispatcher.dispatch(Event{Name: EventNameAppErrorAccept, TargetID: targetIDApp})
			a.dispatcher.dispatch(Event{Name: EventNameAppCmdStop, TargetID: targetIDApp})
			return
		}

		// We only accept the first connection which should be Astilectron, close the next one and stop
		// the app
		if i > 0 {
			astilog.Errorf("Too many TCP connections")
			a.dispatcher.dispatch(Event{Name: EventNameAppTooManyAccept, TargetID: targetIDApp})
			a.dispatcher.dispatch(Event{Name: EventNameAppCmdStop, TargetID: targetIDApp})
			conn.Close()
			return
		}

		// Let the timer know a connection has been accepted
		chanAccepted <- true

		// Create reader and writer
		a.writer = newWriter(conn)
		a.reader = newReader(a.worker.Context(), a.dispatcher, conn)
		go a.reader.read()
	}
}

// execute executes Astilectron in Electron
func (a *Astilectron) execute() (err error) {
	// Log
	astilog.Debug("Executing...")

	// Create command
	var singleInstance string
	if a.options.SingleInstance {
		singleInstance = "true"
	} else {
		singleInstance = "false"
	}
	var cmd = exec.CommandContext(a.worker.Context(), a.paths.AppExecutable(), append([]string{a.paths.AstilectronApplication(), a.listener.Addr().String(), singleInstance}, a.options.ElectronSwitches...)...)
	a.stderrWriter = astikit.NewWriterAdapter(astikit.WriterAdapterOptions{
		Callback: func(i []byte) { astilog.Debugf("Stderr says: %s", i) },
		Split:    []byte("\n"),
	})
	a.stdoutWriter = astikit.NewWriterAdapter(astikit.WriterAdapterOptions{
		Callback: func(i []byte) { astilog.Debugf("Stdout says: %s", i) },
		Split:    []byte("\n"),
	})
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
	var e = synchronousFunc(a.worker.Context(), a, func() {
		err = a.executer(a, cmd)
	}, EventNameAppEventReady)

	// Update display pool
	if e.Displays != nil {
		a.displayPool.update(e.Displays)
	}

	// Create dock
	a.dock = newDock(a.worker.Context(), a.dispatcher, a.identifier, a.writer)

	// Update supported features
	a.supported = e.Supported
	return
}

// watchCmd watches the cmd execution
func (a *Astilectron) watchCmd(cmd *exec.Cmd) {
	// Wait
	cmd.Wait()

	// Check the context to determine whether it was a crash
	if a.worker.Context().Err() != nil {
		astilog.Debug("App has crashed")
		a.dispatcher.dispatch(Event{Name: EventNameAppCrash, TargetID: targetIDApp})
	} else {
		astilog.Debug("App has closed")
		a.dispatcher.dispatch(Event{Name: EventNameAppClose, TargetID: targetIDApp})
	}
	a.dispatcher.dispatch(Event{Name: EventNameAppCmdStop, TargetID: targetIDApp})
}

// Close closes Astilectron properly
func (a *Astilectron) Close() {
	astilog.Debug("Closing...")
	a.worker.Stop()
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
	a.worker.HandleSignals()
}

// Stop orders Astilectron to stop
func (a *Astilectron) Stop() {
	astilog.Debug("Stopping...")
	a.worker.Stop()
}

// Wait is a blocking pattern
func (a *Astilectron) Wait() {
	a.worker.Wait()
}

// Quit quits the app
func (a *Astilectron) Quit() error {
	return a.writer.write(Event{Name: EventNameAppCmdQuit})
}

// Paths returns the paths
func (a *Astilectron) Paths() Paths {
	return *a.paths
}

// Displays returns the displays
func (a *Astilectron) Displays() []*Display {
	return a.displayPool.all()
}

// Dock returns the dock
func (a *Astilectron) Dock() *Dock {
	return a.dock
}

// PrimaryDisplay returns the primary display
func (a *Astilectron) PrimaryDisplay() *Display {
	return a.displayPool.primary()
}

// NewMenu creates a new app menu
func (a *Astilectron) NewMenu(i []*MenuItemOptions) *Menu {
	return newMenu(a.worker.Context(), targetIDApp, i, a.dispatcher, a.identifier, a.writer)
}

// NewWindow creates a new window
func (a *Astilectron) NewWindow(url string, o *WindowOptions) (*Window, error) {
	return newWindow(a.worker.Context(), a.options, a.Paths(), url, o, a.dispatcher, a.identifier, a.writer)
}

// NewWindowInDisplay creates a new window in a specific display
// This overrides the center attribute
func (a *Astilectron) NewWindowInDisplay(d *Display, url string, o *WindowOptions) (*Window, error) {
	if o.X != nil {
		*o.X += d.Bounds().X
	} else {
		o.X = astikit.IntPtr(d.Bounds().X)
	}
	if o.Y != nil {
		*o.Y += d.Bounds().Y
	} else {
		o.Y = astikit.IntPtr(d.Bounds().Y)
	}
	return newWindow(a.worker.Context(), a.options, a.Paths(), url, o, a.dispatcher, a.identifier, a.writer)
}

// NewTray creates a new tray
func (a *Astilectron) NewTray(o *TrayOptions) *Tray {
	return newTray(a.worker.Context(), o, a.dispatcher, a.identifier, a.writer)
}

// NewNotification creates a new notification
func (a *Astilectron) NewNotification(o *NotificationOptions) *Notification {
	return newNotification(a.worker.Context(), o, a.supported != nil && a.supported.Notification != nil && *a.supported.Notification, a.dispatcher, a.identifier, a.writer)
}
