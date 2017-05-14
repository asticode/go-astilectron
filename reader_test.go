package astilectron

import (
	"bytes"
	"io"
	"sync"
	"testing"

	"io/ioutil"

	"github.com/pkg/errors"
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
	var r = newReader(&Dispatcher{}, ioutil.NopCloser(&bytes.Buffer{}))
	assert.True(t, r.isEOFErr(io.EOF))
	assert.True(t, r.isEOFErr(errors.New("read tcp 127.0.0.1:56093->127.0.0.1:56092: wsarecv: An existing connection was forcibly closed by the remote host.")))
	assert.False(t, r.isEOFErr(errors.New("random error")))
}

func TestReader(t *testing.T) {
	// Init
	var mr = &mockedReader{Buffer: bytes.NewBuffer([]byte("{\"name\":\"1\",\"targetId\":\"1\"}\n{\n{\"name\":\"2\",\"targetId\":\"2\"}\n"))}
	var d = newDispatcher()
	defer d.close()
	go d.start()
	var wg = &sync.WaitGroup{}
	var dispatched = []int{}
	d.addListener("1", "1", func(e Event) (deleteListener bool) {
		dispatched = append(dispatched, 1)
		wg.Done()
		return
	})
	d.addListener("2", "2", func(e Event) (deleteListener bool) {
		dispatched = append(dispatched, 2)
		wg.Done()
		return
	})
	wg.Add(2)
	var r = newReader(d, mr)

	// Test read
	go r.read()
	wg.Wait()
	assert.Equal(t, []int{1, 2}, dispatched)

	// Test close
	r.close()
	assert.True(t, mr.c)
}
