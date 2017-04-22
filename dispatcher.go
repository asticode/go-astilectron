package astilectron

import "sync"

// Listener represents a listener executed when an event is dispatched
type Listener func(e Event) (deleteListener bool)

// Listenable represents an object that can listen
type Listenable interface {
	On(eventName string, l Listener)
}

// Dispatcher represents a dispatcher
type Dispatcher struct {
	c  chan Event
	cq chan bool
	l  map[string]map[string][]Listener // Indexed by target ID then by event name
	m  *sync.Mutex
}

// newDispatcher creates a new dispatcher
func newDispatcher() *Dispatcher {
	return &Dispatcher{
		c:  make(chan Event),
		cq: make(chan bool),
		l:  make(map[string]map[string][]Listener),
		m:  &sync.Mutex{},
	}
}

// addListener adds a listener
func (d *Dispatcher) addListener(targetID, eventName string, l Listener) {
	d.m.Lock()
	if _, ok := d.l[targetID]; !ok {
		d.l[targetID] = make(map[string][]Listener)
	}
	d.l[targetID][eventName] = append(d.l[targetID][eventName], l)
	d.m.Unlock()
}

// close closes the dispatcher properly
func (d *Dispatcher) close() {
	if d.cq != nil {
		close(d.cq)
		d.cq = nil
	}
}

// delListener delete a specific listener
func (d *Dispatcher) delListener(targetID, eventName string, index int) {
	d.m.Lock()
	defer d.m.Unlock()
	if _, ok := d.l[targetID]; !ok {
		return
	}
	if _, ok := d.l[targetID][eventName]; !ok {
		return
	}
	d.l[targetID][eventName] = append(d.l[targetID][eventName][:index], d.l[targetID][eventName][index+1:]...)
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
func (d *Dispatcher) listeners(targetID, eventName string) []Listener {
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
