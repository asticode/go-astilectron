package astilectron

import "sync"

// Listener represents a listener executed when an event is dispatched
type Listener func(e Event) (deleteListener bool)

// listenable represents an object that can listen
type listenable interface {
	On(eventName string, l Listener)
}

// dispatcher represents a dispatcher
type dispatcher struct {
	c  chan Event
	cq chan bool
	id int
	// Indexed by target ID then by event name then be listener id
	// We use a map[int]Listener so that deletion is as smooth as possible
	// It means it doesn't store listeners in order
	l map[string]map[string]map[int]Listener
	m *sync.Mutex
}

// newDispatcher creates a new dispatcher
func newDispatcher() *dispatcher {
	return &dispatcher{
		c:  make(chan Event),
		cq: make(chan bool),
		l:  make(map[string]map[string]map[int]Listener),
		m:  &sync.Mutex{},
	}
}

// addListener adds a listener
func (d *dispatcher) addListener(targetID, eventName string, l Listener) {
	d.m.Lock()
	defer d.m.Unlock()
	if _, ok := d.l[targetID]; !ok {
		d.l[targetID] = make(map[string]map[int]Listener)
	}
	if _, ok := d.l[targetID][eventName]; !ok {
		d.l[targetID][eventName] = make(map[int]Listener)
	}
	d.id++
	d.l[targetID][eventName][d.id] = l
}

// close closes the dispatcher properly
func (d *dispatcher) close() {
	if d.cq != nil {
		close(d.cq)
		d.cq = nil
	}
}

// delListener delete a specific listener
func (d *dispatcher) delListener(targetID, eventName string, id int) {
	d.m.Lock()
	defer d.m.Unlock()
	if _, ok := d.l[targetID]; !ok {
		return
	}
	if _, ok := d.l[targetID][eventName]; !ok {
		return
	}
	delete(d.l[targetID][eventName], id)
}

// dispatch dispatches an event
func (d *dispatcher) dispatch(e Event) {
	d.c <- e
}

// start starts the dispatcher and listens to dispatched events
func (d *dispatcher) start() {
	for {
		select {
		case e := <-d.c:
			for id, l := range d.listeners(e.TargetID, e.Name) {
				if l(e) {
					d.delListener(e.TargetID, e.Name, id)
				}
			}
		case <-d.cq:
			return
		}
	}
}

// listeners returns the listeners for a target ID and an event name
func (d *dispatcher) listeners(targetID, eventName string) map[int]Listener {
	d.m.Lock()
	defer d.m.Unlock()
	if _, ok := d.l[targetID]; !ok {
		return map[int]Listener{}
	}
	if _, ok := d.l[targetID][eventName]; !ok {
		return map[int]Listener{}
	}
	return d.l[targetID][eventName]
}
