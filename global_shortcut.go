package astilectron

import (
	"context"
)

const (
	EventNameGlobalShortcutCmdRegister          = "global.shortcut.cmd.register"
	EventNameGlobalShortcutCmdIsRegistered      = "global.shortcut.cmd.is.register"
	EventNameGlobalShortcutCmdUnregister        = "global.shortcut.cmd.unregister"
	EventNameGlobalShortcutCmdUnregisterAll     = "global.shortcut.cmd.unregister.all"
	EventNameGlobalShortcutEventProcessFinished = "global.shortcut.event.process.finished" // Register or Unregister process is finished
	EventNameGlobalShortcutEventTriggered       = "global.shortcut.event.triggered"
)

// GlobalShortcut represents a global shortcut
type GlobalShortcut struct {
	*object
	Accelerator  string `json:"accelerator,omitempty"`  // Accelerator of the global globalShortcuts
	IsRegistered bool   `json:"isRegistered,omitempty"` // Whether the global shortcut is registered
}
type callback func()

var gss = make(map[string]*callback) // Store all registered Global Shortcuts
var obj *object

// InitGlobalShortcuts initializes the globalShortcuts
func InitGlobalShortcuts(ctx context.Context, d *dispatcher, i *identifier, w *writer) {
	obj = newObject(ctx, d, i, w, i.new())
	obj.On(EventNameGlobalShortcutEventTriggered, func(e Event) (deleteListener bool) { // Register the listener for the triggered event
		globalShortcutHandler(e.GlobalShortcut.Accelerator)
		return
	})
}

// GlobalShortcutRegister Register global shortcuts
func GlobalShortcutRegister(accelerator string, callback callback) (isRegistered bool, err error) {

	var gs = GlobalShortcut{Accelerator: accelerator, object: obj}

	// Send an event to astilectron to register the global shortcut
	var event = Event{Name: EventNameGlobalShortcutCmdRegister, TargetID: gs.id, GlobalShortcut: &gs}
	result, err := synchronousEvent(gs.ctx, gs, gs.w, event, EventNameGlobalShortcutEventProcessFinished)

	if err != nil {
		return
	}

	// If registered successfully, add the callback to the map
	if result.GlobalShortcut.IsRegistered {
		gss[accelerator] = &callback
	}

	isRegistered = result.GlobalShortcut.IsRegistered
	return
}

// GlobalShortcutIsRegistered Check if a global shortcut is registered
func GlobalShortcutIsRegistered(accelerator string) (isRegistered bool, err error) {

	var gs = GlobalShortcut{Accelerator: accelerator, object: obj}

	// Send an event to astilectron to check if global shortcut is registered
	var event = Event{Name: EventNameGlobalShortcutCmdIsRegistered, TargetID: gs.id, GlobalShortcut: &gs}
	result, err := synchronousEvent(gs.ctx, gs, gs.w, event, EventNameGlobalShortcutEventProcessFinished)

	if err != nil {
		return
	}

	isRegistered = result.GlobalShortcut.IsRegistered
	return
}

// GlobalShortcutUnregister Unregister a global shortcut
func GlobalShortcutUnregister(accelerator string) (err error) {

	var gs = GlobalShortcut{Accelerator: accelerator, object: obj}

	// Send an event to astilectron to unregister the global shortcut
	var event = Event{Name: EventNameGlobalShortcutCmdUnregister, TargetID: gs.id, GlobalShortcut: &gs}
	_, err = synchronousEvent(gs.ctx, gs, gs.w, event, EventNameGlobalShortcutEventProcessFinished)

	if err != nil {
		return
	}

	// No need to find the callback from the map and delete it
	//  because that event will no longer be triggerred
	// If the same global shortcut is registered again, the original callback will be replaced with the new one

	return
}

// GlobalShortcutUnregisterAll Unregister all global shortcuts
func GlobalShortcutUnregisterAll() (err error) {

	// Send an event to astilectron to unregister all global shortcuts
	var event = Event{Name: EventNameGlobalShortcutCmdUnregisterAll, TargetID: obj.id}
	_, err = synchronousEvent(obj.ctx, obj, obj.w, event, EventNameGlobalShortcutEventProcessFinished)

	if err != nil {
		return
	}

	gss = make(map[string]*callback) // Clear the map

	return
}

// globalShortcutHandler Handle the GlobalShortcut event triggered from astilectron
func globalShortcutHandler(accelerator string) {
	if callback, ok := gss[accelerator]; ok {
		(*callback)()
	}
}
