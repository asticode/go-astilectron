package astilectron

import "sync"

// Listener represents a listener executed when an event is dispatched
type Listener func(payload interface{})

// Dispatcher represents a dispatcher
type Dispatcher struct {
	c  chan Event
	cq chan bool
	l  map[int]map[string][]Listener // Indexed by target ID then by event name
	m  *sync.Mutex
}

// newDispatcher creates a new dispatcher
func newDispatcher() *Dispatcher {
	return &Dispatcher{
		c:  make(chan Event),
		cq: make(chan bool),
		l:  make(map[int]map[string][]Listener),
		m:  &sync.Mutex{},
	}
}

// addListener adds a listener
func (d *Dispatcher) addListener(targetID int, eventName string, l Listener) {
	d.m.Lock()
	if _, ok := d.l[targetID]; !ok {
		d.l[targetID] = make(map[string][]Listener)
	}
	d.l[targetID][eventName] = append(d.l[targetID][eventName], l)
	d.m.Unlock()
}

// close closes the dispatcher properly
func (d *Dispatcher) close() {
	close(d.cq)
}

// Dispatch dispatches an event
func (d *Dispatcher) Dispatch(e Event) {
	d.c <- e
}

// start starts the dispatcher and listens to dispatched events
func (d *Dispatcher) start() {
	for {
		select {
		case e := <-d.c:
			for _, l := range d.listeners(e.TargetID, e.Name) {
				l(e.Payload)
			}
		case <-d.cq:
			return
		}
	}
}

// listeners returns the listeners for a target ID and an event name
func (d *Dispatcher) listeners(targetID int, eventName string) []Listener {
	d.m.Lock()
	defer d.m.Unlock()
	if _, ok := d.l[targetID]; !ok {
		return []Listener{}
	}
	if _, ok := d.l[targetID][eventName]; !ok {
		return []Listener{}
	}
	return d.l[targetID][eventName]
}
