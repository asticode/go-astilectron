package astilectron

import (
	"context"
	"testing"
)

func TestSession_Actions(t *testing.T) {
	// Init
	var d = newDispatcher()
	var i = newIdentifier()
	var wrt = &mockedWriter{}
	var w = newWriter(wrt)
	var s = newSession(context.Background(), d, i, w)

	// Actions
	testObjectAction(t, func() error { return s.ClearCache() }, s.object, wrt, "{\"name\":\"session.cmd.clear.cache\",\"targetID\":\"1\"}\n", EventNameSessionEventClearedCache)
}
