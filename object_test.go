package astilectron

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testObjectAction(t *testing.T, fn func() error, o *object, wrt *mockedWriter, sentEvent, eventNameDone string, receivedEvent *Event) {
	wrt.w = []string{}
	o.cancel()
	err := fn()
	assert.EqualError(t, err, context.Canceled.Error())
	o.ctx, o.cancel = context.WithCancel(context.Background())
	if eventNameDone != "" {
		wrt.fn = func() {
			var event Event
			if receivedEvent != nil {
				event = *receivedEvent
			}
			event.Name = eventNameDone
			event.TargetID = o.id

			o.d.dispatch(event)
		}
	}
	err = fn()
	assert.NoError(t, err)
	assert.Equal(t, []string{sentEvent}, wrt.w)
}
