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
	var w = newWriter(wrt, &logger{})
	var s = newSession(context.Background(), d, i, w)

	// Actions
	testObjectAction(t, func() error { return s.ClearCache() }, s.object, wrt, "{\"name\":\"session.cmd.clear.cache\",\"targetID\":\"1\"}\n", EventNameSessionEventClearedCache, nil)
	testObjectAction(t, func() error { return s.FlushStorage() }, s.object, wrt, "{\"name\":\"session.cmd.flush.storage\",\"targetID\":\"1\"}\n", EventNameSessionEventFlushedStorage, nil)
}
