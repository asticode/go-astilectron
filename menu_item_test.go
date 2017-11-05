package astilectron

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tehsphinx/go-astitools/context"
)

func TestMenuItem_ToEvent(t *testing.T) {
	var o = &MenuItemOptions{Label: PtrStr("1"), SubMenu: []*MenuItemOptions{{Label: PtrStr("2")}, {Label: PtrStr("3")}}}
	var mi = newMenuItem(context.Background(), "main", o, nil, nil, newIdentifier(), nil)
	e := mi.toEvent()
	assert.Equal(t, &EventMenuItem{ID: "1", RootID: "main", Options: o, SubMenu: &EventSubMenu{ID: "2", Items: []*EventMenuItem{{ID: "3", Options: &MenuItemOptions{Label: PtrStr("2")}, RootID: "main"}, {ID: "4", Options: &MenuItemOptions{Label: PtrStr("3")}, RootID: "main"}}, RootID: "main"}}, e)
	assert.Len(t, mi.SubMenu().items, 2)
}

func TestMenuItem_Actions(t *testing.T) {
	// Init
	var c = asticontext.NewCanceller()
	var d = newDispatcher()
	go d.start()
	var i = newIdentifier()
	var wrt = &mockedWriter{}
	var w = newWriter(wrt)
	var mi = newMenuItem(context.Background(), "main", &MenuItemOptions{Label: PtrStr("label")}, c, d, i, w)

	// Actions
	testObjectAction(t, func() error { return mi.SetChecked(true) }, mi.object, wrt, "{\"name\":\""+EventNameMenuItemCmdSetChecked+"\",\"targetID\":\""+mi.id+"\",\"menuItemOptions\":{\"checked\":true}}\n", EventNameMenuItemEventCheckedSet)
	testObjectAction(t, func() error { return mi.SetEnabled(true) }, mi.object, wrt, "{\"name\":\""+EventNameMenuItemCmdSetEnabled+"\",\"targetID\":\""+mi.id+"\",\"menuItemOptions\":{\"enabled\":true}}\n", EventNameMenuItemEventEnabledSet)
	testObjectAction(t, func() error { return mi.SetLabel("test") }, mi.object, wrt, "{\"name\":\""+EventNameMenuItemCmdSetLabel+"\",\"targetID\":\""+mi.id+"\",\"menuItemOptions\":{\"label\":\"test\"}}\n", EventNameMenuItemEventLabelSet)
	testObjectAction(t, func() error { return mi.SetVisible(true) }, mi.object, wrt, "{\"name\":\""+EventNameMenuItemCmdSetVisible+"\",\"targetID\":\""+mi.id+"\",\"menuItemOptions\":{\"visible\":true}}\n", EventNameMenuItemEventVisibleSet)

}
