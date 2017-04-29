package astilectron

import (
	"flag"
	"fmt"
	"net"
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
	versionAstilectron = "0.1.0"
	versionElectron    = "1.6.5"
)

// Vars
var (
	astilectronDirectoryPath = flag.String("astilectron-directory-path", "", "the astilectron directory path")
	boundary                 = []byte("--++__astilectron_boundary__++--")
	validOSes                = []string{"darwin", "linux", "windows"}
)

// Astilectron represents an object capable of interacting with Astilectron
type Astilectron struct {
	canceller   *asticontext.Canceller
	channelQuit chan bool
	dispatcher  *Dispatcher
	identifier  *identifier
	listener    net.Listener
	paths       *Paths
	provisioner Provisioner
	reader      *reader
	writer      *writer
}

// Options represents Astilectron options
type Options struct {
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
		dispatcher:  newDispatcher(),
		identifier:  newIdentifier(),
		provisioner: DefaultProvisioner,
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
	a.dispatcher.Dispatch(Event{Name: EventNameProvisionStart, TargetID: mainTargetID})
	defer a.dispatcher.Dispatch(Event{Name: EventNameProvisionDone, TargetID: mainTargetID})
	return a.provisioner.Provision(a.paths)
}

// listenTCP listens to the first TCP connection coming its way (this should be Astilectron)
func (a *Astilectron) listenTCP() (err error) {
	// Log
	astilog.Debug("Listening...")

	// Listen
	if a.listener, err = net.Listen("tcp", "127.0.0.1:"); err != nil {
		return errors.Wrap(err, "tcp net.Listen failed")
	}
	go func() {
		// We only accept the first connection which should be Astilectron
		var conn net.Conn
		var err error
		if conn, err = a.listener.Accept(); err != nil {
			// TODO Send event and handle it since it's a deal-breaker error
			astilog.Errorf("%s while TCP accepting, not accepting anymore connections", err)
			return
		}

		// Create reader and writer
		a.writer = newWriter(conn)
		a.reader = newReader(a.dispatcher, conn)
		go a.reader.read()
	}()
	return
}

// execute executes Astilectron in Electron
func (a *Astilectron) execute() (err error) {
	// Log
	astilog.Debug("Executing...")

	// Create command
	var ctx, _ = a.canceller.NewContext()
	var cmd = exec.CommandContext(ctx, a.paths.ElectronExecutable(), a.paths.AstilectronApplication(), a.listener.Addr().String())

	// Start command
	synchronousFunc(a, EventNameAppEventReady, func() {
		astilog.Debugf("Starting cmd %s", strings.Join(cmd.Args, " "))
		if err = cmd.Start(); err != nil {
			err = errors.Wrapf(err, "starting cmd %s failed", strings.Join(cmd.Args, " "))
			return
		}
	})
	return
}

// Close closes Astilectron properly
func (a *Astilectron) Close() {
	astilog.Debug("Closing...")
	a.canceller.Cancel()
	a.dispatcher.close()
	a.listener.Close()
	a.reader.close()
	a.writer.close()
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
	if a.channelQuit != nil {
		close(a.channelQuit)
		a.channelQuit = nil
	}
}

// Wait is a blocking pattern
func (a *Astilectron) Wait() {
	for {
		select {
		case <-a.channelQuit:
			return
		}
	}
}
