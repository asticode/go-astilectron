package astilectron

import (
	"context"
	"sync"
)

const (
	EventNameGlobalShortcutsCmdRegister          = "global.shortcuts.cmd.register"
	EventNameGlobalShortcutsCmdIsRegistered      = "global.shortcuts.cmd.is.register"
	EventNameGlobalShortcutsCmdUnregister        = "global.shortcuts.cmd.unregister"
	EventNameGlobalShortcutsCmdUnregisterAll     = "global.shortcuts.cmd.unregister.all"
	EventNameGlobalShortcutsEventRegistered      = "global.shortcuts.event.registered"
	EventNameGlobalShortcutsEventIsRegistered    = "global.shortcuts.event.is.registered"
	EventNameGlobalShortcutsEventUnregistered    = "global.shortcuts.event.unregistered"
	EventNameGlobalShortcutsEventUnregisteredAll = "global.shortcuts.event.unregistered.all"
	EventNameGlobalShortcutEventTriggered        = "global.shortcuts.event.triggered"
)

type globalShortcutsCallback func()

// GlobalShortcuts represents a global shortcut
type GlobalShortcuts struct {
	*object
	m         *sync.Mutex
	callbacks map[string]*globalShortcutsCallback // Store all registered Global Shortcuts
}

func newGlobalShortcuts(ctx context.Context, d *dispatcher, i *identifier, w *writer) (gs *GlobalShortcuts) {

	gs = &GlobalShortcuts{object: newObject(ctx, d, i, w, i.new()), m: new(sync.Mutex), callbacks: make(map[string]*globalShortcutsCallback)}
	gs.On(EventNameGlobalShortcutEventTriggered, func(e Event) (deleteListener bool) { // Register the listener for the triggered event
		gs.execCallback(e.GlobalShortcuts.Accelerator)
		return
	})
	return
}

// Register Register global shortcuts
func (gs *GlobalShortcuts) Register(accelerator string, callback globalShortcutsCallback) (isRegistered bool, err error) {

	if err = gs.ctx.Err(); err != nil {
		return
	}

	// Send an event to astilectron to register the global shortcut
	var event = Event{Name: EventNameGlobalShortcutsCmdRegister, TargetID: gs.id, GlobalShortcuts: &EventGlobalShortcuts{Accelerator: accelerator}}
	result, err := synchronousEvent(gs.ctx, gs, gs.w, event, EventNameGlobalShortcutsEventRegistered)

	if err != nil {
		return
	}

	// If registered successfully, add the callback to the map
	if result.GlobalShortcuts.IsRegistered {
		gs.m.Lock()
		gs.callbacks[accelerator] = &callback
		gs.m.Unlock()
	}

	isRegistered = result.GlobalShortcuts.IsRegistered
	return
}

// IsRegistered Check if a global shortcut is registered
func (gs *GlobalShortcuts) IsRegistered(accelerator string) (isRegistered bool, err error) {

	if err = gs.ctx.Err(); err != nil {
		return
	}

	// Send an event to astilectron to check if global shortcut is registered
	var event = Event{Name: EventNameGlobalShortcutsCmdIsRegistered, TargetID: gs.id, GlobalShortcuts: &EventGlobalShortcuts{Accelerator: accelerator}}
	result, err := synchronousEvent(gs.ctx, gs, gs.w, event, EventNameGlobalShortcutsEventIsRegistered)

	if err != nil {
		return
	}

	isRegistered = result.GlobalShortcuts.IsRegistered
	return
}

// Unregister Unregister a global shortcut
func (gs *GlobalShortcuts) Unregister(accelerator string) (err error) {

	if err = gs.ctx.Err(); err != nil {
		return
	}

	// Send an event to astilectron to unregister the global shortcut
	var event = Event{Name: EventNameGlobalShortcutsCmdUnregister, TargetID: gs.id, GlobalShortcuts: &EventGlobalShortcuts{Accelerator: accelerator}}
	_, err = synchronousEvent(gs.ctx, gs, gs.w, event, EventNameGlobalShortcutsEventUnregistered)

	if err != nil {
		return
	}

	// No need to find the callback from the map and delete it
	//  because that event will no longer be triggerred
	// If the same global shortcut is registered again, the original callback will be replaced with the new one

	return
}

// UnregisterAll Unregister all global shortcuts
func (gs *GlobalShortcuts) UnregisterAll() (err error) {

	if err = gs.ctx.Err(); err != nil {
		return
	}

	// Send an event to astilectron to unregister all global shortcuts
	var event = Event{Name: EventNameGlobalShortcutsCmdUnregisterAll, TargetID: gs.id}
	_, err = synchronousEvent(gs.ctx, gs, gs.w, event, EventNameGlobalShortcutsEventUnregisteredAll)

	if err != nil {
		return
	}

	gs.m.Lock()
	gs.callbacks = make(map[string]*globalShortcutsCallback) // Clear the map
	gs.m.Unlock()

	return
}

// execCallback Execute the GlobalShortcuts event triggered from astilectron
func (gs *GlobalShortcuts) execCallback(accelerator string) {
	gs.m.Lock()
	if callback, ok := gs.callbacks[accelerator]; ok {
		(*callback)()
	}
	gs.m.Unlock()
}
