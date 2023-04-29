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
	var w = newWriter(wrt, &logger{})
	var tr = newTray(context.Background(), &TrayOptions{
		Image:   astikit.StrPtr("/path/to/image"),
		Tooltip: astikit.StrPtr("tooltip"),
	}, d, i, w)
	var m = tr.NewMenu([]*MenuItemOptions{
		{
			Label: astikit.StrPtr("Root 1"),
			SubMenu: []*MenuItemOptions{
				{Label: astikit.StrPtr("Item 1")},
				{Label: astikit.StrPtr("Item 2")},
				{Type: MenuItemTypeSeparator},
				{Label: astikit.StrPtr("Item 3")},
			},
		}})
	var p = PositionOptions{
		X: astikit.IntPtr(250),
		Y: astikit.IntPtr(250),
	}

	// Actions
	testObjectAction(t, func() error { return tr.Create() }, tr.object, wrt, "{\"name\":\""+EventNameTrayCmdCreate+"\",\"targetID\":\""+tr.id+"\",\"trayOptions\":{\"image\":\"/path/to/image\",\"tooltip\":\"tooltip\"}}\n", EventNameTrayEventCreated, nil)
	testObjectAction(t, func() error { return tr.SetImage("test") }, tr.object, wrt, "{\"name\":\""+EventNameTrayCmdSetImage+"\",\"targetID\":\""+tr.id+"\",\"image\":\"test\"}\n", EventNameTrayEventImageSet, nil)
	testObjectAction(t, func() error { return tr.PopUpContextMenu( &TrayPopUpOptions{}) }, tr.object, wrt, "{\"name\":\""+EventNameTrayCmdPopUpContextMenu+"\",\"targetID\":\""+tr.id+"\",\"menuPopupOptions\":{}}\n", EventNameTrayEventPoppedUpContextMenu, nil)
	testObjectAction(t, func() error { return tr.PopUpContextMenu( &TrayPopUpOptions{Position: p}) }, tr.object, wrt, "{\"name\":\""+EventNameTrayCmdPopUpContextMenu+"\",\"targetID\":\""+tr.id+"\",\"menuPopupOptions\":{\"x\":250,\"y\":250}}\n", EventNameTrayEventPoppedUpContextMenu, nil)
	testObjectAction(t, func() error { return tr.PopUpContextMenu( &TrayPopUpOptions{Menu: m}) }, tr.object, wrt, "{\"name\":\""+EventNameTrayCmdPopUpContextMenu+"\",\"targetID\":\""+tr.id+"\",\"menu\":{\"id\":\""+m.id+"\",\"items\":[{\"id\":\"3\",\"options\":{\"label\":\"Root 1\"},\"rootId\":\""+m.rootID+"\",\"submenu\":{\"id\":\"4\",\"items\":[{\"id\":\"5\",\"options\":{\"label\":\"Item 1\"},\"rootId\":\""+m.rootID+"\"},{\"id\":\"6\",\"options\":{\"label\":\"Item 2\"},\"rootId\":\""+m.rootID+"\"},{\"id\":\"7\",\"options\":{\"type\":\"separator\"},\"rootId\":\""+m.rootID+"\"},{\"id\":\"8\",\"options\":{\"label\":\"Item 3\"},\"rootId\":\""+m.rootID+"\"}],\"rootId\":\""+m.rootID+"\"}}],\"rootId\":\""+m.rootID+"\"},\"menuPopupOptions\":{}}\n", EventNameTrayEventPoppedUpContextMenu, nil)
	testObjectAction(t, func() error { return tr.PopUpContextMenu( &TrayPopUpOptions{Menu: m, Position: p}) }, tr.object, wrt, "{\"name\":\""+EventNameTrayCmdPopUpContextMenu+"\",\"targetID\":\""+tr.id+"\",\"menu\":{\"id\":\""+m.id+"\",\"items\":[{\"id\":\"3\",\"options\":{\"label\":\"Root 1\"},\"rootId\":\""+m.rootID+"\",\"submenu\":{\"id\":\"4\",\"items\":[{\"id\":\"5\",\"options\":{\"label\":\"Item 1\"},\"rootId\":\""+m.rootID+"\"},{\"id\":\"6\",\"options\":{\"label\":\"Item 2\"},\"rootId\":\""+m.rootID+"\"},{\"id\":\"7\",\"options\":{\"type\":\"separator\"},\"rootId\":\""+m.rootID+"\"},{\"id\":\"8\",\"options\":{\"label\":\"Item 3\"},\"rootId\":\""+m.rootID+"\"}],\"rootId\":\""+m.rootID+"\"}}],\"rootId\":\""+m.rootID+"\"},\"menuPopupOptions\":{\"x\":250,\"y\":250}}\n", EventNameTrayEventPoppedUpContextMenu, nil)
	testObjectAction(t, func() error { return tr.Destroy() }, tr.object, wrt, "{\"name\":\""+EventNameTrayCmdDestroy+"\",\"targetID\":\""+tr.id+"\"}\n", EventNameTrayEventDestroyed, nil)
	assert.True(t, tr.ctx.Err() != nil)
}

func TestTray_NewMenu(t *testing.T) {
	a, err := New(nil, Options{})
	assert.NoError(t, err)
	tr := a.NewTray(&TrayOptions{})
	m := tr.NewMenu([]*MenuItemOptions{})
	assert.Equal(t, tr.id, m.rootID)
}
