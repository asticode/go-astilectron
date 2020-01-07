package astilectron

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockedReader represents a mocked reader
type mockedReader struct {
	*bytes.Buffer
	c bool
}

// Close implements the io.Close interface
func (r *mockedReader) Close() error {
	r.c = true
	return nil
}

func TestReader_IsEOFErr(t *testing.T) {
	var r = newReader(context.Background(), &logger{}, &dispatcher{}, ioutil.NopCloser(&bytes.Buffer{}))
	assert.True(t, r.isEOFErr(io.EOF))
	assert.True(t, r.isEOFErr(errors.New("read tcp 127.0.0.1:56093->127.0.0.1:56092: wsarecv: An existing connection was forcibly closed by the remote host.")))
	assert.False(t, r.isEOFErr(errors.New("random error")))
}

func TestReader(t *testing.T) {
	// Init
	var mr = &mockedReader{Buffer: bytes.NewBuffer([]byte("{\"name\":\"1\",\"targetId\":\"1\"}\n{\n{\"name\":\"2\",\"targetId\":\"2\"}\n"))}
	var d = newDispatcher()
	var wg = &sync.WaitGroup{}
	var dispatched = []int{}
	var dispatchedMutex = sync.Mutex{}
	d.addListener("1", "1", func(e Event) (deleteListener bool) {
		dispatchedMutex.Lock()
		dispatched = append(dispatched, 1)
		dispatchedMutex.Unlock()
		wg.Done()
		return
	})
	d.addListener("2", "2", func(e Event) (deleteListener bool) {
		dispatchedMutex.Lock()
		dispatched = append(dispatched, 2)
		dispatchedMutex.Unlock()
		wg.Done()
		return
	})
	wg.Add(2)
	var r = newReader(context.Background(), &logger{}, d, mr)

	// Test read
	go r.read()
	wg.Wait()
	assert.Contains(t, dispatched, 1)
	assert.Contains(t, dispatched, 2)

	// Test close
	r.close()
	assert.True(t, mr.c)
}
