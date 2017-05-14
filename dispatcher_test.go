package astilectron

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDispatcher(t *testing.T) {
	// Init
	var d = newDispatcher()
	defer d.close()
	go d.start()
	var wg = &sync.WaitGroup{}
	var dispatched = []int{}

	// Test adding listener
	d.addListener("1", "1", func(e Event) (deleteListener bool) {
		dispatched = append(dispatched, 1)
		wg.Done()
		return
	})
	d.addListener("1", "1", func(e Event) (deleteListener bool) {
		dispatched = append(dispatched, 2)
		wg.Done()
		return true
	})
	d.addListener("1", "1", func(e Event) (deleteListener bool) {
		dispatched = append(dispatched, 3)
		wg.Done()
		return true
	})
	d.addListener("1", "2", func(e Event) (deleteListener bool) {
		dispatched = append(dispatched, 4)
		wg.Done()
		return
	})
	assert.Len(t, d.l["1"]["1"], 3)

	// Test dispatch
	wg.Add(4)
	d.Dispatch(Event{Name: "2", TargetID: "1"})
	d.Dispatch(Event{Name: "1", TargetID: "1"})
	wg.Wait()
	for _, v := range []int{1, 2, 3, 4} {
		assert.Contains(t, dispatched, v)
	}
	assert.Len(t, d.l["1"]["1"], 1)

	// Test close
	d.close()
	assert.Nil(t, d.cq)
}
