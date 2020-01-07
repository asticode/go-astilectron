package astilectron

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/asticode/go-astikit"
	"github.com/stretchr/testify/assert"
)

// mockedHandler is a mocked handler
type mockedHandler struct {
	e bool
}

func (h *mockedHandler) readFile(rw http.ResponseWriter, path string) {
	var b, err = ioutil.ReadFile(path)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(b)
}

// ServeHTTP implements the http.Handler interface
func (h *mockedHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if h.e {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	switch r.URL.Path {
	case "/provisioner/astilectron":
		h.readFile(rw, "testdata/provisioner/astilectron/astilectron.zip")
	case "/provisioner/electron/darwin":
		h.readFile(rw, "testdata/provisioner/electron/darwin/electron.zip")
	case "/provisioner/electron/linux":
		h.readFile(rw, "testdata/provisioner/electron/linux/electron.zip")
	case "/provisioner/electron/windows":
		h.readFile(rw, "testdata/provisioner/electron/windows/electron.zip")
	default:
		rw.Write([]byte("body"))
	}
}

var tempPathCount int

func mockedTempPath() string {
	tempPathCount++
	return fmt.Sprintf("testdata/tmp/%d", tempPathCount)
}

func TestDownload(t *testing.T) {
	// Init
	var mh = &mockedHandler{e: true}
	var s = httptest.NewServer(mh)
	var dst = mockedTempPath()
	var d = astikit.NewHTTPDownloader(astikit.HTTPDownloaderOptions{})

	// Test failed download
	err := Download(context.Background(), &logger{}, d, s.URL, dst)
	assert.Error(t, err)
	_, err = os.Stat(dst)
	assert.True(t, os.IsNotExist(err))

	// Test successful download
	mh.e = false
	err = Download(context.Background(), &logger{}, d, s.URL, dst)
	assert.NoError(t, err)
	defer os.Remove(dst)
	b, err := ioutil.ReadFile(dst)
	assert.NoError(t, err)
	assert.Equal(t, "body", string(b))
}

// mockedDisembedder is a mocked disembedder
func mockedDisembedder(src string) ([]byte, error) {
	switch src {
	case "astilectron":
		return ioutil.ReadFile("testdata/provisioner/astilectron/astilectron.zip")
	case "electron/linux":
		return ioutil.ReadFile("testdata/provisioner/electron/linux/electron.zip")
	case "test":
		return []byte("body"), nil
	default:
		return []byte{}, errors.New("invalid")
	}
}

func TestDisembed(t *testing.T) {
	// Init
	var dst = mockedTempPath()

	// Test failed disembed
	err := Disembed(context.Background(), &logger{}, mockedDisembedder, "invalid", dst)
	assert.EqualError(t, err, "disembedding invalid failed: invalid")

	// Test successful disembed
	err = Disembed(context.Background(), &logger{}, mockedDisembedder, "test", dst)
	assert.NoError(t, err)
	defer os.Remove(dst)
	b, err := ioutil.ReadFile(dst)
	assert.NoError(t, err)
	assert.Equal(t, "body", string(b))
}

func TestPtr(t *testing.T) {
	assert.Equal(t, true, *astikit.BoolPtr(true))
	assert.Equal(t, 1, *astikit.IntPtr(1))
	assert.Equal(t, "1", *astikit.StrPtr("1"))
}

// mockedListenable is a mocked listenable
type mockedListenable struct {
	d  *dispatcher
	id string
}

// On implements the listenable interface
func (m *mockedListenable) On(eventName string, l Listener) {
	m.d.addListener(m.id, eventName, l)
}

func TestSynchronousFunc(t *testing.T) {
	// Init
	var d = newDispatcher()
	var l = &mockedListenable{d: d, id: "1"}
	var done bool
	var m sync.Mutex
	l.On("done", func(e Event) bool {
		m.Lock()
		defer m.Unlock()
		done = true
		return false
	})

	// Test canceller cancel
	ctx, cancel := context.WithCancel(context.Background())
	var _ = synchronousFunc(ctx, l, func() { cancel() }, "done")
	assert.False(t, done)

	// Test done event
	var ed = Event{Name: "done", TargetID: "1"}
	var e = synchronousFunc(context.Background(), l, func() { d.dispatch(ed) }, "done")
	m.Lock()
	assert.True(t, done)
	m.Unlock()
	assert.Equal(t, ed, e)
}

func TestSynchronousEvent(t *testing.T) {
	// Init
	var d = newDispatcher()
	var ed = Event{Name: "done", TargetID: "1"}
	var mw = &mockedWriter{fn: func() { d.dispatch(ed) }}
	var w = newWriter(mw, &logger{})
	var l = &mockedListenable{d: d, id: "1"}
	var done bool
	var m sync.Mutex
	l.On("done", func(e Event) bool {
		m.Lock()
		defer m.Unlock()
		done = true
		return false
	})
	var ei = Event{Name: "order", TargetID: "1"}

	// Test successful synchronous event
	var e, err = synchronousEvent(context.Background(), l, w, ei, "done")
	assert.NoError(t, err)
	m.Lock()
	assert.True(t, done)
	m.Unlock()
	assert.Equal(t, ed, e)
	assert.Equal(t, []string{"{\"name\":\"order\",\"targetID\":\"1\"}\n"}, mw.w)
}
