package astilectron

import (
	"testing"

	"github.com/asticode/go-astitools/context"
	"github.com/stretchr/testify/assert"
)

func TestObject_IsActionable(t *testing.T) {
	// Init
	var c = asticontext.NewCanceller()
	var o = newObject(nil, c, nil, newIdentifier(), nil)

	// Test success
	assert.NoError(t, o.isActionable())

	// Test object destroyed
	o.cancel()
	assert.EqualError(t, o.isActionable(), ErrObjectDestroyed.Error())

	// Test canceller cancelled
	c.Cancel()
	assert.EqualError(t, o.isActionable(), ErrCancellerCancelled.Error())
}

func testObjectAction(t *testing.T, fn func() error, o *object, wrt *mockedWriter, sentEvent, eventNameDone string) {
	wrt.w = []string{}
	o.c.Cancel()
	err := fn()
	assert.EqualError(t, err, ErrCancellerCancelled.Error())
	o.c = asticontext.NewCanceller()
	o.ctx, o.cancel = o.c.NewContext()
	if eventNameDone != "" {
		wrt.fn = func() { o.d.dispatch(Event{Name: eventNameDone, TargetID: o.id}) }
	}
	err = fn()
	assert.NoError(t, err)
	assert.Equal(t, []string{sentEvent}, wrt.w)
}
