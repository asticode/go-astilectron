package astilectron

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tehsphinx/go-astitools/context"
)

func TestMenu_ToEvent(t *testing.T) {
	var m = newMenu(nil, "main", []*MenuItemOptions{{Label: PtrStr("1")}, {Label: PtrStr("2")}}, asticontext.NewCanceller(), newDispatcher(), newIdentifier(), nil)
	e := m.toEvent()
	assert.Equal(t, &EventMenu{EventSubMenu: &EventSubMenu{ID: "1", Items: []*EventMenuItem{{ID: "2", Options: &MenuItemOptions{Label: PtrStr("1")}, RootID: "main"}, {ID: "3", Options: &MenuItemOptions{Label: PtrStr("2")}, RootID: "main"}}, RootID: "main"}}, e)
}

func TestMenu_Actions(t *testing.T) {
	// Init
	var c = asticontext.NewCanceller()
	var d = newDispatcher()
	go d.start()
	var i = newIdentifier()
	var wrt = &mockedWriter{}
	var w = newWriter(wrt)
	var m = newMenu(context.Background(), "main", []*MenuItemOptions{{Label: PtrStr("1")}, {Label: PtrStr("2")}}, c, d, i, w)

	// Actions
	testObjectAction(t, func() error { return m.Create() }, m.object, wrt, "{\"name\":\""+EventNameMenuCmdCreate+"\",\"targetID\":\""+m.id+"\",\"menu\":{\"id\":\"1\",\"items\":[{\"id\":\"2\",\"options\":{\"label\":\"1\"},\"rootId\":\"main\"},{\"id\":\"3\",\"options\":{\"label\":\"2\"},\"rootId\":\"main\"}],\"rootId\":\"main\"}}\n", EventNameMenuEventCreated)
	testObjectAction(t, func() error { return m.Destroy() }, m.object, wrt, "{\"name\":\""+EventNameMenuCmdDestroy+"\",\"targetID\":\""+m.id+"\",\"menu\":{\"id\":\"1\",\"items\":[{\"id\":\"2\",\"options\":{\"label\":\"1\"},\"rootId\":\"main\"},{\"id\":\"3\",\"options\":{\"label\":\"2\"},\"rootId\":\"main\"}],\"rootId\":\"main\"}}\n", EventNameMenuEventDestroyed)
	assert.True(t, m.IsDestroyed())
}
