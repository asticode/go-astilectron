package astilectron

import (
	"context"
	"testing"

	"github.com/asticode/go-astikit"
	"github.com/stretchr/testify/assert"
)

func TestMenu_ToEvent(t *testing.T) {
	var m = newMenu(context.Background(), targetIDApp, []*MenuItemOptions{{Label: astikit.StrPtr("1")}, {Label: astikit.StrPtr("2")}}, newDispatcher(), newIdentifier(), nil)
	e := m.toEvent()
	assert.Equal(t, &EventMenu{EventSubMenu: &EventSubMenu{ID: "1", Items: []*EventMenuItem{{ID: "2", Options: &MenuItemOptions{Label: astikit.StrPtr("1")}, RootID: targetIDApp}, {ID: "3", Options: &MenuItemOptions{Label: astikit.StrPtr("2")}, RootID: targetIDApp}}, RootID: targetIDApp}}, e)
}

func TestMenu_Actions(t *testing.T) {
	// Init
	var d = newDispatcher()
	var i = newIdentifier()
	var wrt = &mockedWriter{}
	var w = newWriter(wrt, &logger{})
	var m = newMenu(context.Background(), targetIDApp, []*MenuItemOptions{{Label: astikit.StrPtr("1")}, {Label: astikit.StrPtr("2")}}, d, i, w)

	// Actions
	testObjectAction(t, func() error { return m.Create() }, m.object, wrt, "{\"name\":\""+EventNameMenuCmdCreate+"\",\"targetID\":\""+m.id+"\",\"menu\":{\"id\":\"1\",\"items\":[{\"id\":\"2\",\"options\":{\"label\":\"1\"},\"rootId\":\""+targetIDApp+"\"},{\"id\":\"3\",\"options\":{\"label\":\"2\"},\"rootId\":\""+targetIDApp+"\"}],\"rootId\":\""+targetIDApp+"\"}}\n", EventNameMenuEventCreated, nil)
	testObjectAction(t, func() error { return m.Destroy() }, m.object, wrt, "{\"name\":\""+EventNameMenuCmdDestroy+"\",\"targetID\":\""+m.id+"\",\"menu\":{\"id\":\"1\",\"items\":[{\"id\":\"2\",\"options\":{\"label\":\"1\"},\"rootId\":\""+targetIDApp+"\"},{\"id\":\"3\",\"options\":{\"label\":\"2\"},\"rootId\":\""+targetIDApp+"\"}],\"rootId\":\""+targetIDApp+"\"}}\n", EventNameMenuEventDestroyed, nil)
	assert.True(t, m.ctx.Err() != nil)
}
