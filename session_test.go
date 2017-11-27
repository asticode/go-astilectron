package astilectron

import (
	"context"
	"testing"

	"github.com/asticode/go-astitools/context"
)

func TestSession_Actions(t *testing.T) {
	// Init
	var c = asticontext.NewCanceller()
	var d = newDispatcher()
	var i = newIdentifier()
	var wrt = &mockedWriter{}
	var w = newWriter(wrt)
	var s = newSession(context.Background(), c, d, i, w)

	// Actions
	testObjectAction(t, func() error { return s.ClearCache() }, s.object, wrt, "{\"name\":\"session.cmd.clear.cache\",\"targetID\":\"1\"}\n", EventNameSessionEventClearedCache)
}
