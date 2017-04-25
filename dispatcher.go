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
	l  map[string]map[string][]Listener // Indexed by target ID then by event name
	m  *sync.Mutex
}

// newDispatcher creates a new dispatcher
func newDispatcher() *dispatcher {
	return &dispatcher{
		c:  make(chan Event),
		cq: make(chan bool),
		l:  make(map[string]map[string][]Listener),
		m:  &sync.Mutex{},
	}
}

// addListener adds a listener
func (d *dispatcher) addListener(targetID, eventName string, l Listener) {
	d.m.Lock()
	if _, ok := d.l[targetID]; !ok {
		d.l[targetID] = make(map[string][]Listener)
	}
	d.l[targetID][eventName] = append(d.l[targetID][eventName], l)
	d.m.Unlock()
}

// close closes the dispatcher properly
func (d *dispatcher) close() {
	if d.cq != nil {
		close(d.cq)
		d.cq = nil
	}
}

// delListener delete a specific listener
func (d *dispatcher) delListener(targetID, eventName string, index int) {
	d.m.Lock()
	defer d.m.Unlock()
	if _, ok := d.l[targetID]; !ok {
		return
	}
	if _, ok := d.l[targetID][eventName]; !ok {
		return
	}
	if len(d.l[targetID][eventName]) <= 1 {
		d.l[targetID][eventName] = []Listener{}
		return
	}
	d.l[targetID][eventName] = append(d.l[targetID][eventName][:index], d.l[targetID][eventName][index+1:]...)
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
			for i, l := range d.listeners(e.TargetID, e.Name) {
				if deleteListener := l(e); deleteListener {
					d.delListener(e.TargetID, e.Name, i)
				}
			}
		case <-d.cq:
			return
		}
	}
}

// listeners returns the listeners for a target ID and an event name
func (d *dispatcher) listeners(targetID, eventName string) []Listener {
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
