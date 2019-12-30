package astilectron

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testObjectAction(t *testing.T, fn func() error, o *object, wrt *mockedWriter, sentEvent, eventNameDone string) {
	wrt.w = []string{}
	o.cancel()
	err := fn()
	assert.EqualError(t, err, context.Canceled.Error())
	o.ctx, o.cancel = context.WithCancel(context.Background())
	if eventNameDone != "" {
		wrt.fn = func() { o.d.dispatch(Event{Name: eventNameDone, TargetID: o.id}) }
	}
	err = fn()
	assert.NoError(t, err)
	assert.Equal(t, []string{sentEvent}, wrt.w)
}
