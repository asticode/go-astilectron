package astilectron

import (
	"net"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAstilectron_Provision(t *testing.T) {
	// Init
	var o = Options{BaseDirectoryPath: mockedTempPath()}
	defer os.RemoveAll(o.BaseDirectoryPath)
	a, err := New(o)
	assert.NoError(t, err)
	defer a.dispatcher.close()
	go a.dispatcher.start()
	a.SetProvisioner(NewDisembedderProvisioner(mockedDisembedder, "astilectron", "electron/linux"))
	var hasStarted, hasStopped bool
	a.On(EventNameProvisionAstilectronMoved, func(e Event) bool {
		hasStarted = true
		return false
	})
	var wg = &sync.WaitGroup{}
	a.On(EventNameProvisionElectronFinished, func(e Event) bool {
		hasStopped = true
		wg.Done()
		return false
	})
	wg.Add(1)

	// Test provision is successful and sends the correct events
	err = a.provision()
	assert.NoError(t, err)
	wg.Wait()
	assert.True(t, hasStarted)
	assert.True(t, hasStopped)
}

func TestAstilectron_WatchNoAccept(t *testing.T) {
	// Init
	a, err := New(Options{})
	assert.NoError(t, err)
	defer a.dispatcher.close()
	go a.dispatcher.start()
	var isStopped bool
	var wg = &sync.WaitGroup{}
	a.On(EventNameAppCmdStop, func(e Event) bool {
		isStopped = true
		wg.Done()
		return false
	})
	c := make(chan bool)

	// Test success
	go func() {
		time.Sleep(50 * time.Microsecond)
		c <- true
	}()
	a.watchNoAccept(time.Second, c)
	assert.False(t, isStopped)

	// Test failure
	wg.Add(1)
	a.watchNoAccept(time.Nanosecond, c)
	wg.Wait()
	assert.True(t, isStopped)
}

// mockedListener implements the net.Listener interface
type mockedListener struct {
	c chan bool
	e chan bool
}

func (l mockedListener) Accept() (net.Conn, error) {
	for {
		select {
		case <-l.c:
			return mockedConn{}, nil
		case <-l.e:
			return nil, errors.New("invalid")
		}
	}
}
func (l mockedListener) Close() error   { return nil }
func (l mockedListener) Addr() net.Addr { return nil }

// mockedConn implements the net.Conn interface
type mockedConn struct{}

func (c mockedConn) Read(b []byte) (n int, err error)   { return }
func (c mockedConn) Write(b []byte) (n int, err error)  { return }
func (c mockedConn) Close() error                       { return nil }
func (c mockedConn) LocalAddr() net.Addr                { return nil }
func (c mockedConn) RemoteAddr() net.Addr               { return nil }
func (c mockedConn) SetDeadline(t time.Time) error      { return nil }
func (c mockedConn) SetReadDeadline(t time.Time) error  { return nil }
func (c mockedConn) SetWriteDeadline(t time.Time) error { return nil }

// mockedAddr implements the net.Addr interface
type mockedAddr struct{}

func (a mockedAddr) Network() string { return "" }
func (a mockedAddr) String() string  { return "" }

func TestAstilectron_AcceptTCP(t *testing.T) {
	// Init
	a, err := New(Options{})
	assert.NoError(t, err)
	defer a.Close()
	go a.dispatcher.start()
	var l = &mockedListener{c: make(chan bool), e: make(chan bool)}
	a.listener = l
	var isStopped bool
	var wg = &sync.WaitGroup{}
	a.On(EventNameAppCmdStop, func(e Event) bool {
		isStopped = true
		wg.Done()
		return false
	})
	c := make(chan bool)
	var isAccepted bool
	go func() {
		for {
			select {
			case <-c:
				isAccepted = true
				wg.Done()
			}
		}
	}()
	go a.acceptTCP(c)

	// Test accepted
	wg.Add(1)
	l.c <- true
	wg.Wait()
	assert.True(t, isAccepted)
	assert.False(t, isStopped)

	// Test refused
	isAccepted = false
	wg.Add(1)
	l.c <- true
	wg.Wait()
	assert.False(t, isAccepted)
	assert.True(t, isStopped)

	// Test error accept
	go a.acceptTCP(c)
	isStopped = false
	wg.Add(1)
	l.e <- true
	wg.Wait()
	assert.False(t, isAccepted)
	assert.True(t, isStopped)
}

func TestAstilectron_ExecuteCmd(t *testing.T) {
	// Init
	a, err := New(Options{})
	assert.NoError(t, err)
	defer a.Close()
	go a.dispatcher.start()
	var wg = &sync.WaitGroup{}

	// Test success
	var cmd = exec.Command(os.Getenv("GOROOT")+"go", "version")
	wg.Add(1)
	go func() {
		a.executeCmd(cmd)
		wg.Done()
	}()
	a.dispatcher.Dispatch(Event{Name: EventNameAppEventReady, Displays: &EventDisplays{All: []*DisplayOptions{{ID: PtrInt(1)}}, Primary: &DisplayOptions{ID: PtrInt(1)}}, TargetID: mainTargetID})
	wg.Wait()
	assert.Len(t, a.Displays(), 1)
	assert.Equal(t, 1, *a.PrimaryDisplay().o.ID)
}

func TestIsValidOS(t *testing.T) {
	assert.NoError(t, validateOS("darwin"))
	assert.NoError(t, validateOS("linux"))
	assert.NoError(t, validateOS("windows"))
	assert.Error(t, validateOS("invalid"))
}

func TestAstilectron_Wait(t *testing.T) {
	a, err := New(Options{})
	assert.NoError(t, err)
	a.HandleSignals()
	go func() {
		time.Sleep(20 * time.Microsecond)
		p, err := os.FindProcess(os.Getpid())
		assert.NoError(t, err)
		p.Signal(os.Interrupt)
	}()
	a.Wait()
}
