package astilectron

import (
	"net"
	"os"
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
	a.SetProvisioner(NewDisembedderProvisioner(mockedDisembedder, "astilectron", "electron/linux"))

	// Test provision is successful
	err = a.provision()
	assert.NoError(t, err)
}

func TestAstilectron_WatchNoAccept(t *testing.T) {
	// Init
	a, err := New(Options{})
	assert.NoError(t, err)
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
				return
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

func TestIsValidOS(t *testing.T) {
	assert.True(t, IsValidOS("darwin"))
	assert.True(t, IsValidOS("linux"))
	assert.True(t, IsValidOS("windows"))
	assert.False(t, IsValidOS("invalid"))
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

func TestAstilectron_NewMenu(t *testing.T) {
	a, err := New(Options{})
	assert.NoError(t, err)
	m := a.NewMenu([]*MenuItemOptions{})
	assert.Equal(t, targetIDApp, m.rootID)
}

func TestAstilectron_Actions(t *testing.T) {
	// Init
	a, err := New(Options{})
	assert.NoError(t, err)
	defer a.Close()
	wrt := &mockedWriter{}
	a.writer = newWriter(wrt)

	// Actions
	err = a.Quit()
	assert.NoError(t, err)
	assert.Equal(t, []string{"{\"name\":\"app.cmd.quit\"}\n"}, wrt.w)
}
