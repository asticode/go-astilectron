package astilectron

import (
	"net/url"

	"github.com/asticode/go-astitools/context"
	"github.com/asticode/go-astitools/url"
	"github.com/pkg/errors"
)

// Window event names
const (
	EventNameWindowCmdBlur                     = "window.cmd.blur"
	EventNameWindowCmdCenter                   = "window.cmd.center"
	EventNameWindowCmdClose                    = "window.cmd.close"
	EventNameWindowCmdCreate                   = "window.cmd.create"
	EventNameWindowCmdDestroy                  = "window.cmd.destroy"
	EventNameWindowCmdFocus                    = "window.cmd.focus"
	EventNameWindowCmdHide                     = "window.cmd.hide"
	EventNameWindowCmdMaximize                 = "window.cmd.maximize"
	EventNameWindowCmdMessage                  = "window.cmd.message"
	EventNameWindowCmdMinimize                 = "window.cmd.minimize"
	EventNameWindowCmdMove                     = "window.cmd.move"
	EventNameWindowCmdResize                   = "window.cmd.resize"
	EventNameWindowCmdRestore                  = "window.cmd.restore"
	EventNameWindowCmdShow                     = "window.cmd.show"
	EventNameWindowCmdUnmaximize               = "window.cmd.unmaximize"
	EventNameWindowCmdWebContentsCloseDevTools = "window.cmd.web.contents.close.dev.tools"
	EventNameWindowCmdWebContentsOpenDevTools  = "window.cmd.web.contents.open.dev.tools"
	EventNameWindowEventBlur                   = "window.event.blur"
	EventNameWindowEventClosed                 = "window.event.closed"
	EventNameWindowEventDidFinishLoad          = "window.event.did.finish.load"
	EventNameWindowEventFocus                  = "window.event.focus"
	EventNameWindowEventHide                   = "window.event.hide"
	EventNameWindowEventMaximize               = "window.event.maximize"
	EventNameWindowEventMessage                = "window.event.message"
	EventNameWindowEventMinimize               = "window.event.minimize"
	EventNameWindowEventMove                   = "window.event.move"
	EventNameWindowEventReadyToShow            = "window.event.ready.to.show"
	EventNameWindowEventResize                 = "window.event.resize"
	EventNameWindowEventRestore                = "window.event.restore"
	EventNameWindowEventShow                   = "window.event.show"
	EventNameWindowEventUnmaximize             = "window.event.unmaximize"
	EventNameWindowEventUnresponsive           = "window.event.unresponsive"
)

// Title bar styles
var (
	TitleBarStyleDefault     = PtrStr("default")
	TitleBarStyleHidden      = PtrStr("hidden")
	TitleBarStyleHiddenInset = PtrStr("hidden-inset")
)

// Window represents a window
// TODO Add missing window options
// TODO Add missing window methods
// TODO Add missing window events
type Window struct {
	*object
	o   *WindowOptions
	url *url.URL
}

// WindowOptions represents window options
// We must use pointers since GO doesn't handle optional fields whereas NodeJS does. Use PtrBool, PtrInt or PtrStr
// to fill the struct
// https://github.com/electron/electron/blob/v1.6.5/docs/api/browser-window.md
type WindowOptions struct {
	AcceptFirstMouse       *bool   `json:"acceptFirstMouse,omitempty"`
	AlwaysOnTop            *bool   `json:"alwaysOnTop,omitempty"`
	AutoHideMenuBar        *bool   `json:"autoHideMenuBar,omitempty"`
	BackgroundColor        *string `json:"backgroundColor,omitempty"`
	Center                 *bool   `json:"center,omitempty"`
	Closable               *bool   `json:"closable,omitempty"`
	DisableAutoHideCursor  *bool   `json:"disableAutoHideCursor,omitempty"`
	EnableLargerThanScreen *bool   `json:"enableLargerThanScreen,omitempty"`
	Focusable              *bool   `json:"focusable,omitempty"`
	Frame                  *bool   `json:"frame,omitempty"`
	Fullscreen             *bool   `json:"fullscreen,omitempty"`
	Fullscreenable         *bool   `json:"fullscreenable,omitempty"`
	HasShadow              *bool   `json:"hasShadow,omitempty"`
	Height                 *int    `json:"height,omitempty"`
	Icon                   *string `json:"icon,omitempty"`
	Kiosk                  *bool   `json:"kiosk,omitempty"`
	MaxHeight              *int    `json:"maxHeight,omitempty"`
	Maximizable            *bool   `json:"maximizable,omitempty"`
	MaxWidth               *int    `json:"maxWidth,omitempty"`
	MinHeight              *int    `json:"minHeight,omitempty"`
	Minimizable            *bool   `json:"minimizable,omitempty"`
	MinWidth               *int    `json:"minWidth,omitempty"`
	Modal                  *bool   `json:"modal,omitempty"`
	Movable                *bool   `json:"movable,omitempty"`
	Resizable              *bool   `json:"resizable,omitempty"`
	Show                   *bool   `json:"show,omitempty"`
	SkipTaskbar            *bool   `json:"skipTaskbar,omitempty"`
	Title                  *string `json:"title,omitempty"`
	TitleBarStyle          *string `json:"titleBarStyle,omitempty"`
	Transparent            *bool   `json:"transparent,omitempty"`
	UseContentSize         *bool   `json:"useContentSize,omitempty"`
	Width                  *int    `json:"width,omitempty"`
	X                      *int    `json:"x,omitempty"`
	Y                      *int    `json:"y,omitempty"`
}

// newWindow creates a new window
func newWindow(o Options, url string, wo *WindowOptions, c *asticontext.Canceller, d *Dispatcher, i *identifier, wrt *writer) (w *Window, err error) {
	// Init
	w = &Window{
		o:      wo,
		object: newObject(nil, c, d, i, wrt),
	}

	// Check app details
	if wo.Icon == nil && o.AppIconDefaultPath != "" {
		wo.Icon = PtrStr(o.AppIconDefaultPath)
	}
	if wo.Title == nil && o.AppName != "" {
		wo.Title = PtrStr(o.AppName)
	}

	// Make sure the window's context is cancelled once the closed event is received
	w.On(EventNameWindowEventClosed, func(e Event) (deleteListener bool) {
		w.cancel()
		return true
	})

	// Parse url
	if w.url, err = astiurl.Parse(url); err != nil {
		err = errors.Wrapf(err, "parsing url %s failed", url)
		return
	}
	return
}

// NewMenu creates a new window menu
func (w *Window) NewMenu(i []*MenuItemOptions) *Menu {
	return newMenu(w.ctx, w, i, w.c, w.d, w.i, w.w)
}

// Blur blurs the window
func (w *Window) Blur() (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(w.c, w, w.w, Event{Name: EventNameWindowCmdBlur, TargetID: w.id}, EventNameWindowEventBlur)
	return
}

// Center centers the window
func (w *Window) Center() (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(w.c, w, w.w, Event{Name: EventNameWindowCmdCenter, TargetID: w.id}, EventNameWindowEventMove)
	return
}

// Close closes the window
func (w *Window) Close() (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(w.c, w, w.w, Event{Name: EventNameWindowCmdClose, TargetID: w.id}, EventNameWindowEventClosed)
	return
}

// CloseDevTools closes the dev tools
func (w *Window) CloseDevTools() (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	return w.w.write(Event{Name: EventNameWindowCmdWebContentsCloseDevTools, TargetID: w.id})
}

// Create creates the window
// We wait for EventNameWindowEventDidFinishLoad since we need the web content to be fully loaded before being able to
// send messages to it
func (w *Window) Create() (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(w.c, w, w.w, Event{Name: EventNameWindowCmdCreate, TargetID: w.id, URL: w.url.String(), WindowOptions: w.o}, EventNameWindowEventDidFinishLoad)
	return
}

// Destroy destroys the window
func (w *Window) Destroy() (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(w.c, w, w.w, Event{Name: EventNameWindowCmdDestroy, TargetID: w.id}, EventNameWindowEventClosed)
	return
}

// Focus focuses on the window
func (w *Window) Focus() (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(w.c, w, w.w, Event{Name: EventNameWindowCmdFocus, TargetID: w.id}, EventNameWindowEventFocus)
	return
}

// Hide hides the window
func (w *Window) Hide() (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(w.c, w, w.w, Event{Name: EventNameWindowCmdHide, TargetID: w.id}, EventNameWindowEventHide)
	return
}

// OpenDevTools opens the dev tools
func (w *Window) OpenDevTools() (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	return w.w.write(Event{Name: EventNameWindowCmdWebContentsOpenDevTools, TargetID: w.id})
}

// Maximize maximizes the window
func (w *Window) Maximize() (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(w.c, w, w.w, Event{Name: EventNameWindowCmdMaximize, TargetID: w.id}, EventNameWindowEventMaximize)
	return
}

// Minimize minimizes the window
func (w *Window) Minimize() (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(w.c, w, w.w, Event{Name: EventNameWindowCmdMinimize, TargetID: w.id}, EventNameWindowEventMinimize)
	return
}

// Move moves the window
func (w *Window) Move(x, y int) (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	w.o.X = PtrInt(x)
	w.o.Y = PtrInt(y)
	_, err = synchronousEvent(w.c, w, w.w, Event{Name: EventNameWindowCmdMove, TargetID: w.id, WindowOptions: &WindowOptions{X: w.o.X, Y: w.o.Y}}, EventNameWindowEventMove)
	return
}

// MoveInDisplay moves the window in the proper display
func (w *Window) MoveInDisplay(d *Display, x, y int) error {
	return w.Move(d.Bounds().X+x, d.Bounds().Y+y)
}

// Resize resizes the window
func (w *Window) Resize(width, height int) (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	w.o.Height = PtrInt(height)
	w.o.Width = PtrInt(width)
	_, err = synchronousEvent(w.c, w, w.w, Event{Name: EventNameWindowCmdResize, TargetID: w.id, WindowOptions: &WindowOptions{Height: w.o.Height, Width: w.o.Width}}, EventNameWindowEventResize)
	return
}

// Restore restores the window
func (w *Window) Restore() (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(w.c, w, w.w, Event{Name: EventNameWindowCmdRestore, TargetID: w.id}, EventNameWindowEventRestore)
	return
}

// Send sends a message to the inner JS of the Web content of the window
func (w *Window) Send(message interface{}) (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	return w.w.write(Event{Message: newEventMessage(message), Name: EventNameWindowCmdMessage, TargetID: w.id})
}

// Show shows the window
func (w *Window) Show() (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(w.c, w, w.w, Event{Name: EventNameWindowCmdShow, TargetID: w.id}, EventNameWindowEventShow)
	return
}

// Unmaximize unmaximize the window
func (w *Window) Unmaximize() (err error) {
	if err = w.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(w.c, w, w.w, Event{Name: EventNameWindowCmdUnmaximize, TargetID: w.id}, EventNameWindowEventUnmaximize)
	return
}
