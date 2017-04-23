package astilectron

import (
	"net/url"

	"github.com/pkg/errors"
)

// Window represents a window
type Window struct {
	d   *Dispatcher
	id  string
	o   *WindowOptions
	url *url.URL
	w   *writer
}

// WindowOptions represents window options
// We must use pointers since GO doesn't handle optional fields whereas NodeJS does. Use PtrBool, PtrInt or PtrStr
// to fill the struct
// https://github.com/electron/electron/blob/master/docs/api/browser-window.md#new-browserwindowoptions
// TODO Add missing attributes
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
	Transparent            *bool   `json:"transparent,omitempty"`
	UseContentSize         *bool   `json:"useContentSize,omitempty"`
	Width                  *int    `json:"width,omitempty"`
	X                      *int    `json:"x,omitempty"`
	Y                      *int    `json:"y,omitempty"`
}

// NewWindow creates a new window
func (a *Astilectron) NewWindow(url string, o *WindowOptions) (w *Window, err error) {
	// Init
	w = &Window{
		d:  a.dispatcher,
		id: a.identifier.new(),
		o:  o,
		w:  a.writer,
	}

	// Parse url
	if w.url, err = parseURL(url); err != nil {
		err = errors.Wrapf(err, "parsing url %s failed", url)
		return
	}
	return
}

// On implements the Listenable interface
func (w *Window) On(eventName string, l Listener) {
	w.d.addListener(w.id, eventName, l)
}

// Create creates the window
func (w *Window) Create() error {
	return synchronousEvent(w, w.w, Event{Name: EventNameWindowCreate, TargetID: w.id, URL: w.url.String(), WindowOptions: w.o}, EventNameWindowCreateDone)
}

// Show shows the window
func (w *Window) Show() error {
	return synchronousEvent(w, w.w, Event{Name: EventNameWindowShow, TargetID: w.id}, EventNameWindowShowDone)
}
