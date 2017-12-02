package astilectron

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDispatcher(t *testing.T) {
	// Init
	var d = newDispatcher()
	var wg = sync.WaitGroup{}
	var dispatched = []int{}
	var m sync.Mutex

	// Test adding listener
	d.addListener("1", "1", func(e Event) (deleteListener bool) {
		m.Lock()
		dispatched = append(dispatched, 1)
		m.Unlock()
		wg.Done()
		return
	})
	d.addListener("1", "1", func(e Event) (deleteListener bool) {
		m.Lock()
		dispatched = append(dispatched, 2)
		m.Unlock()
		wg.Done()
		return true
	})
	d.addListener("1", "1", func(e Event) (deleteListener bool) {
		m.Lock()
		dispatched = append(dispatched, 3)
		m.Unlock()
		wg.Done()
		return true
	})
	d.addListener("1", "2", func(e Event) (deleteListener bool) {
		m.Lock()
		dispatched = append(dispatched, 4)
		m.Unlock()
		wg.Done()
		return
	})
	assert.Len(t, d.l["1"]["1"], 3)

	// Test dispatch
	wg.Add(4)
	d.dispatch(Event{Name: "2", TargetID: "1"})
	d.dispatch(Event{Name: "1", TargetID: "1"})
	wg.Wait()
	for _, v := range []int{1, 2, 3, 4} {
		assert.Contains(t, dispatched, v)
	}
	assert.Len(t, d.listeners("1", "1"), 1)
}
