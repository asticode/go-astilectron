package astilectron

import (
	"context"
	"testing"

	"github.com/asticode/go-astitools/context"
	"github.com/stretchr/testify/assert"
)

func TestSubMenu_ToEvent(t *testing.T) {
	// App sub menu
	var s = newSubMenu(nil, "main", []*MenuItemOptions{{Label: PtrStr("1")}, {Label: PtrStr("2")}}, asticontext.NewCanceller(), newDispatcher(), newIdentifier(), nil)
	e := s.toEvent()
	assert.Equal(t, &EventSubMenu{ID: "1", Items: []*EventMenuItem{{ID: "2", Options: &MenuItemOptions{Label: PtrStr("1")}, RootID: "main"}, {ID: "3", Options: &MenuItemOptions{Label: PtrStr("2")}, RootID: "main"}}, RootID: "main"}, e)

	// Window sub menu
	var i = newIdentifier()
	w, err := newWindow(Options{}, "http://test.com", &WindowOptions{}, asticontext.NewCanceller(), newDispatcher(), i, nil)
	assert.NoError(t, err)
	s = newSubMenu(nil, w.id, []*MenuItemOptions{{Label: PtrStr("1")}, {Label: PtrStr("2")}}, asticontext.NewCanceller(), newDispatcher(), i, nil)
	e = s.toEvent()
	assert.Equal(t, &EventSubMenu{ID: "2", Items: []*EventMenuItem{{ID: "3", Options: &MenuItemOptions{Label: PtrStr("1")}, RootID: "1"}, {ID: "4", Options: &MenuItemOptions{Label: PtrStr("2")}, RootID: "1"}}, RootID: "1"}, e)
}

func TestSubMenu_SubMenu(t *testing.T) {
	var o = []*MenuItemOptions{
		{},
		{SubMenu: []*MenuItemOptions{
			{},
			{SubMenu: []*MenuItemOptions{
				{},
				{},
				{},
			}},
			{},
		}},
		{},
	}
	var m = newMenu(context.Background(), "main", o, nil, newDispatcher(), newIdentifier(), nil)
	s, err := m.SubMenu(0, 1)
	assert.EqualError(t, err, "No submenu at 0")
	s, err = m.SubMenu(1)
	assert.NoError(t, err)
	assert.Len(t, s.items, 3)
	s, err = m.SubMenu(1, 0)
	assert.EqualError(t, err, "No submenu at 1:0")
	s, err = m.SubMenu(1, 1)
	assert.NoError(t, err)
	assert.Len(t, s.items, 3)
	s, err = m.SubMenu(1, 3)
	assert.EqualError(t, err, "Submenu at 1 has 3 items, invalid index 3")
}

func TestSubMenu_Item(t *testing.T) {
	var o = []*MenuItemOptions{
		{Label: PtrStr("1")},
		{Label: PtrStr("2"), SubMenu: []*MenuItemOptions{
			{Label: PtrStr("2-1")},
			{Label: PtrStr("2-2"), SubMenu: []*MenuItemOptions{
				{Label: PtrStr("2-2-1")},
				{Label: PtrStr("2-2-2")},
				{Label: PtrStr("2-2-3")},
			}},
			{Label: PtrStr("2-3")},
		}},
		{Label: PtrStr("3")},
	}
	var m = newMenu(context.Background(), "main", o, nil, newDispatcher(), newIdentifier(), nil)
	i, err := m.Item(3)
	assert.EqualError(t, err, "Submenu has 3 items, invalid index 3")
	i, err = m.Item(0)
	assert.NoError(t, err)
	assert.Equal(t, "1", *i.o.Label)
	i, err = m.Item(1, 3)
	assert.EqualError(t, err, "Submenu has 3 items, invalid index 3")
	i, err = m.Item(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, "2-3", *i.o.Label)
	i, err = m.Item(1, 1, 0)
	assert.NoError(t, err)
	assert.Equal(t, "2-2-1", *i.o.Label)
}

func TestSubMenu_Actions(t *testing.T) {
	// Init
	var c = asticontext.NewCanceller()
	var d = newDispatcher()
	go d.start()
	var i = newIdentifier()
	var wrt = &mockedWriter{}
	var w = newWriter(wrt)
	var s = newSubMenu(nil, "main", []*MenuItemOptions{{Label: PtrStr("0")}}, c, d, i, w)

	// Actions
	var mi = s.NewItem(&MenuItemOptions{Label: PtrStr("1")})
	testObjectAction(t, func() error { return s.Append(mi) }, s.object, wrt, "{\"name\":\""+EventNameSubMenuCmdAppend+"\",\"targetID\":\""+s.id+"\",\"menuItem\":{\"id\":\"3\",\"options\":{\"label\":\"1\"},\"rootId\":\"main\"}}\n", EventNameSubMenuEventAppended)
	assert.Len(t, s.items, 2)
	assert.Equal(t, "1", *s.items[1].o.Label)
	mi = s.NewItem(&MenuItemOptions{Label: PtrStr("2")})
	err := s.Insert(3, mi)
	assert.EqualError(t, err, "Submenu has 2 items, position 3 is invalid")
	testObjectAction(t, func() error { return s.Insert(1, mi) }, s.object, wrt, "{\"name\":\""+EventNameSubMenuCmdInsert+"\",\"targetID\":\""+s.id+"\",\"menuItem\":{\"id\":\"4\",\"options\":{\"label\":\"2\"},\"rootId\":\"main\"},\"menuItemPosition\":1}\n", EventNameSubMenuEventInserted)
	assert.Len(t, s.items, 3)
	assert.Equal(t, "2", *s.items[1].o.Label)
	testObjectAction(t, func() error {
		return s.Popup(&MenuPopupOptions{PositionOptions: PositionOptions{X: PtrInt(1), Y: PtrInt(2)}})
	}, s.object, wrt, "{\"name\":\""+EventNameSubMenuCmdPopup+"\",\"targetID\":\""+s.id+"\",\"menuPopupOptions\":{\"x\":1,\"y\":2}}\n", EventNameSubMenuEventPoppedUp)
	testObjectAction(t, func() error {
		return s.PopupInWindow(&Window{object: &object{id: "2"}}, &MenuPopupOptions{PositionOptions: PositionOptions{X: PtrInt(1), Y: PtrInt(2)}})
	}, s.object, wrt, "{\"name\":\""+EventNameSubMenuCmdPopup+"\",\"targetID\":\""+s.id+"\",\"menuPopupOptions\":{\"x\":1,\"y\":2},\"windowId\":\"2\"}\n", EventNameSubMenuEventPoppedUp)
	testObjectAction(t, func() error { return s.ClosePopup() }, s.object, wrt, "{\"name\":\""+EventNameSubMenuCmdClosePopup+"\",\"targetID\":\""+s.id+"\"}\n", EventNameSubMenuEventClosedPopup)
	testObjectAction(t, func() error { return s.ClosePopupInWindow(&Window{object: &object{id: "2"}}) }, s.object, wrt, "{\"name\":\""+EventNameSubMenuCmdClosePopup+"\",\"targetID\":\""+s.id+"\",\"windowId\":\"2\"}\n", EventNameSubMenuEventClosedPopup)
}
