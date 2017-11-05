package astilectron

import (
	"context"
	"errors"

	"github.com/tehsphinx/go-astitools/context"
)

// Object errors
var (
	ErrCancellerCancelled = errors.New("canceller.cancelled")
	ErrObjectDestroyed    = errors.New("object.destroyed")
)

// object represents a base object
type object struct {
	cancel context.CancelFunc
	ctx    context.Context
	c      *asticontext.Canceller
	d      *Dispatcher
	i      *identifier
	id     string
	w      *writer
}

// newObject returns a new base object
func newObject(parentCtx context.Context, c *asticontext.Canceller, d *Dispatcher, i *identifier, w *writer) (o *object) {
	o = &object{
		c:  c,
		d:  d,
		i:  i,
		id: i.new(),
		w:  w,
	}
	if parentCtx != nil {
		o.ctx, o.cancel = context.WithCancel(parentCtx)
	} else {
		o.ctx, o.cancel = c.NewContext()
	}
	return
}

// isActionable checks whether any type of action is allowed on the window
func (o *object) isActionable() error {
	if o.c.Cancelled() {
		return ErrCancellerCancelled
	} else if o.IsDestroyed() {
		return ErrObjectDestroyed
	}
	return nil
}

// IsDestroyed checks whether the window has been destroyed
func (o *object) IsDestroyed() bool {
	return o.ctx.Err() != nil
}

// On implements the Listenable interface
func (o *object) On(eventName string, l Listener) {
	o.d.addListener(o.id, eventName, l)
}
