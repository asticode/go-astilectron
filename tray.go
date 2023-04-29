package astilectron

import (
	"context"

	"github.com/asticode/go-astikit"
)

// Tray event names
const (
	EventNameTrayCmdCreate                = "tray.cmd.create"
	EventNameTrayCmdDestroy               = "tray.cmd.destroy"
	EventNameTrayCmdSetImage              = "tray.cmd.set.image"
	EventNameTrayCmdPopUpContextMenu      = "tray.cmd.popup.contextmenu"
	EventNameTrayEventClicked             = "tray.event.clicked"
	EventNameTrayEventCreated             = "tray.event.created"
	EventNameTrayEventDestroyed           = "tray.event.destroyed"
	EventNameTrayEventDoubleClicked       = "tray.event.double.clicked"
	EventNameTrayEventImageSet            = "tray.event.image.set"
	EventNameTrayEventRightClicked        = "tray.event.right.clicked"
	EventNameTrayEventPoppedUpContextMenu = "tray.event.poppedup.contextmenu"
)

// Tray represents a tray
type Tray struct {
	*object
	o *TrayOptions
}

// TrayOptions represents tray options
// We must use pointers since GO doesn't handle optional fields whereas NodeJS does. Use astikit.BoolPtr, astikit.IntPtr or astikit.StrPtr
// to fill the struct
// https://github.com/electron/electron/blob/v1.8.1/docs/api/tray.md
type TrayOptions struct {
	Image   *string `json:"image,omitempty"`
	Tooltip *string `json:"tooltip,omitempty"`
}

// TrayPopUpOptions represents Tray PopUpContextMenu options
type TrayPopUpOptions struct {
	Menu     *Menu
	Position PositionOptions
}

// newTray creates a new tray
func newTray(ctx context.Context, o *TrayOptions, d *dispatcher, i *identifier, wrt *writer) (t *Tray) {
	// Init
	t = &Tray{
		o:      o,
		object: newObject(ctx, d, i, wrt, i.new()),
	}

	// Make sure the tray's context is cancelled once the destroyed event is received
	t.On(EventNameTrayEventDestroyed, func(e Event) (deleteListener bool) {
		t.cancel()
		return true
	})
	return
}

// Create creates the tray
func (t *Tray) Create() (err error) {
	if err = t.ctx.Err(); err != nil {
		return
	}
	var e = Event{Name: EventNameTrayCmdCreate, TargetID: t.id, TrayOptions: t.o}
	_, err = synchronousEvent(t.ctx, t, t.w, e, EventNameTrayEventCreated)
	return
}

// Destroy destroys the tray
func (t *Tray) Destroy() (err error) {
	if err = t.ctx.Err(); err != nil {
		return
	}
	_, err = synchronousEvent(t.ctx, t, t.w, Event{Name: EventNameTrayCmdDestroy, TargetID: t.id}, EventNameTrayEventDestroyed)
	return
}

// NewMenu creates a new tray menu
func (t *Tray) NewMenu(i []*MenuItemOptions) *Menu {
	return newMenu(t.ctx, t.id, i, t.d, t.i, t.w)
}

// SetImage sets the tray image
func (t *Tray) SetImage(image string) (err error) {
	if err = t.ctx.Err(); err != nil {
		return
	}
	t.o.Image = astikit.StrPtr(image)
	_, err = synchronousEvent(t.ctx, t, t.w, Event{Name: EventNameTrayCmdSetImage, Image: image, TargetID: t.id}, EventNameTrayEventImageSet)
	return
}

// PopUpContextMenu pops up the context menu of the tray icon.
// When menu is passed, the menu will be shown instead of the tray icon's context menu.
// The position is only available on Windows, and it is (0, 0) by default.
// https://www.electronjs.org/docs/latest/api/tray#traypopupcontextmenumenu-position-macos-windows
func (t *Tray) PopUpContextMenu(p *TrayPopUpOptions) (err error) {
	var em *EventMenu
	if err = t.ctx.Err(); err != nil {
		return
	}
	if p.Menu != nil {
		em = p.Menu.toEvent()
	}
	var e = Event{Name: EventNameTrayCmdPopUpContextMenu, TargetID: t.id, Menu: em, MenuPopupOptions: &MenuPopupOptions{PositionOptions: p.Position}}
	_, err = synchronousEvent(t.ctx, t, t.w, e, EventNameTrayEventPoppedUpContextMenu)
	return
}