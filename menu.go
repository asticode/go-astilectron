package astilectron

import (
	"context"

	"github.com/asticode/go-astitools/context"
)

// Menu event names
const (
	EventNameMenuCmdCreate    = "menu.cmd.create"
	EventNameMenuEventCreated = "menu.event.created"
)

// Menu represents a menu
// https://github.com/electron/electron/blob/v1.6.5/docs/api/menu.md
type Menu struct {
	*subMenu
}

// newMenu creates a new menu
func newMenu(ctx context.Context, root interface{}, items []*MenuItemOptions, c *asticontext.Canceller, d *dispatcher, i *identifier, w *writer) *Menu {
	return &Menu{newSubMenu(ctx, root, items, c, d, i, w)}
}

// toEvent returns the menu in the proper event format
func (m *Menu) toEvent() *EventMenu {
	return &EventMenu{m.subMenu.toEvent()}
}

// Create creates the menu
func (m *Menu) Create() (err error) {
	if err = m.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(m.c, m, m.w, Event{Name: EventNameMenuCmdCreate, TargetID: m.id, Menu: m.toEvent()}, EventNameMenuEventCreated)
	return
}
