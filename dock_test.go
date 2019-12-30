package astilectron

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDock_Actions(t *testing.T) {
	// Init
	var d = newDispatcher()
	var i = newIdentifier()
	var wrt = &mockedWriter{}
	var w = newWriter(wrt)
	var dck = newDock(context.Background(), d, i, w)

	// Actions
	testObjectAction(t, func() error {
		_, err := dck.Bounce(DockBounceTypeCritical)
		return err
	}, dck.object, wrt, "{\"name\":\""+eventNameDockCmdBounce+"\",\"targetID\":\""+dck.id+"\",\"bounceType\":\"critical\"}\n", eventNameDockEventBouncing)
	testObjectAction(t, func() error { return dck.BounceDownloads("/path/to/file") }, dck.object, wrt, "{\"name\":\""+eventNameDockCmdBounceDownloads+"\",\"targetID\":\""+dck.id+"\",\"filePath\":\"/path/to/file\"}\n", eventNameDockEventDownloadsBouncing)
	testObjectAction(t, func() error { return dck.CancelBounce(1) }, dck.object, wrt, "{\"name\":\""+eventNameDockCmdCancelBounce+"\",\"targetID\":\""+dck.id+"\",\"id\":1}\n", eventNameDockEventBouncingCancelled)
	testObjectAction(t, func() error { return dck.Hide() }, dck.object, wrt, "{\"name\":\""+eventNameDockCmdHide+"\",\"targetID\":\""+dck.id+"\"}\n", eventNameDockEventHidden)
	testObjectAction(t, func() error { return dck.SetBadge("badge") }, dck.object, wrt, "{\"name\":\""+eventNameDockCmdSetBadge+"\",\"targetID\":\""+dck.id+"\",\"badge\":\"badge\"}\n", eventNameDockEventBadgeSet)
	testObjectAction(t, func() error { return dck.SetIcon("/path/to/icon") }, dck.object, wrt, "{\"name\":\""+eventNameDockCmdSetIcon+"\",\"targetID\":\""+dck.id+"\",\"image\":\"/path/to/icon\"}\n", eventNameDockEventIconSet)
	testObjectAction(t, func() error { return dck.Show() }, dck.object, wrt, "{\"name\":\""+eventNameDockCmdShow+"\",\"targetID\":\""+dck.id+"\"}\n", eventNameDockEventShown)
}

func TestDock_NewMenu(t *testing.T) {
	var d = newDispatcher()
	var i = newIdentifier()
	var wrt = &mockedWriter{}
	var w = newWriter(wrt)
	var dck = newDock(context.Background(), d, i, w)
	m := dck.NewMenu([]*MenuItemOptions{})
	assert.Equal(t, dck.id, m.rootID)
}
