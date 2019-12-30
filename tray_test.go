package astilectron

import (
	"context"
	"testing"

	"github.com/asticode/go-astikit"
	"github.com/stretchr/testify/assert"
)

func TestTray_Actions(t *testing.T) {
	// Init
	var d = newDispatcher()
	var i = newIdentifier()
	var wrt = &mockedWriter{}
	var w = newWriter(wrt)
	var tr = newTray(context.Background(), &TrayOptions{
		Image:   astikit.StrPtr("/path/to/image"),
		Tooltip: astikit.StrPtr("tooltip"),
	}, d, i, w)

	// Actions
	testObjectAction(t, func() error { return tr.Create() }, tr.object, wrt, "{\"name\":\""+EventNameTrayCmdCreate+"\",\"targetID\":\""+tr.id+"\",\"trayOptions\":{\"image\":\"/path/to/image\",\"tooltip\":\"tooltip\"}}\n", EventNameTrayEventCreated)
	testObjectAction(t, func() error { return tr.SetImage("test") }, tr.object, wrt, "{\"name\":\""+EventNameTrayCmdSetImage+"\",\"targetID\":\""+tr.id+"\",\"image\":\"test\"}\n", EventNameTrayEventImageSet)
	testObjectAction(t, func() error { return tr.Destroy() }, tr.object, wrt, "{\"name\":\""+EventNameTrayCmdDestroy+"\",\"targetID\":\""+tr.id+"\"}\n", EventNameTrayEventDestroyed)
	assert.True(t, tr.ctx.Err() != nil)
}

func TestTray_NewMenu(t *testing.T) {
	a, err := New(Options{})
	assert.NoError(t, err)
	tr := a.NewTray(&TrayOptions{})
	m := tr.NewMenu([]*MenuItemOptions{})
	assert.Equal(t, tr.id, m.rootID)
}
