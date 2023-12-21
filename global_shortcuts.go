package astilectron

import (
	"context"
	"sync"
)

const (
	EventNameGlobalShortcutsCmdRegister          = "global.shortcuts.cmd.register"
	EventNameGlobalShortcutsCmdIsRegistered      = "global.shortcuts.cmd.is.registered"
	EventNameGlobalShortcutsCmdUnregister        = "global.shortcuts.cmd.unregister"
	EventNameGlobalShortcutsCmdUnregisterAll     = "global.shortcuts.cmd.unregister.all"
	EventNameGlobalShortcutsEventRegistered      = "global.shortcuts.event.registered"
	EventNameGlobalShortcutsEventIsRegistered    = "global.shortcuts.event.is.registered"
	EventNameGlobalShortcutsEventUnregistered    = "global.shortcuts.event.unregistered"
	EventNameGlobalShortcutsEventUnregisteredAll = "global.shortcuts.event.unregistered.all"
	EventNameGlobalShortcutEventTriggered        = "global.shortcuts.event.triggered"
)

type globalShortcutsCallback func()

// GlobalShortcuts represents global shortcuts
type GlobalShortcuts struct {
	*object
	m         *sync.Mutex
	callbacks map[string]globalShortcutsCallback
}

func newGlobalShortcuts(ctx context.Context, d *dispatcher, i *identifier, w *writer) (gs *GlobalShortcuts) {
	gs = &GlobalShortcuts{
		object:    newObject(ctx, d, i, w, i.new()),
		m:         new(sync.Mutex),
		callbacks: make(map[string]globalShortcutsCallback),
	}
	gs.On(EventNameGlobalShortcutEventTriggered, func(e Event) (deleteListener bool) { // Register the listener for the triggered event
		gs.m.Lock()
		callback, ok := gs.callbacks[e.GlobalShortcuts.Accelerator]
		gs.m.Unlock()
		if ok {
			(callback)()
		}
		return
	})
	return
}

// Register registers a global shortcut
func (gs *GlobalShortcuts) Register(accelerator string, callback globalShortcutsCallback) (isRegistered bool, err error) {
	if err = gs.ctx.Err(); err != nil {
		return
	}

	// Send an event to astilectron to register the global shortcut
	result, err := synchronousEvent(gs.ctx, gs, gs.w, Event{Name: EventNameGlobalShortcutsCmdRegister, TargetID: gs.id, GlobalShortcuts: &EventGlobalShortcuts{Accelerator: accelerator}}, EventNameGlobalShortcutsEventRegistered)
	if err != nil {
		return
	}

	// If registered successfully, add the callback to the map
	if result.GlobalShortcuts != nil {
		if result.GlobalShortcuts.IsRegistered {
			gs.m.Lock()
			gs.callbacks[accelerator] = callback
			gs.m.Unlock()
		}
		isRegistered = result.GlobalShortcuts.IsRegistered
	}
	return
}

// IsRegistered checks whether a global shortcut is registered
func (gs *GlobalShortcuts) IsRegistered(accelerator string) (isRegistered bool, err error) {
	if err = gs.ctx.Err(); err != nil {
		return
	}

	// Send an event to astilectron to check if global shortcut is registered
	result, err := synchronousEvent(gs.ctx, gs, gs.w, Event{Name: EventNameGlobalShortcutsCmdIsRegistered, TargetID: gs.id, GlobalShortcuts: &EventGlobalShortcuts{Accelerator: accelerator}}, EventNameGlobalShortcutsEventIsRegistered)
	if err != nil {
		return
	}

	if result.GlobalShortcuts != nil {
		isRegistered = result.GlobalShortcuts.IsRegistered
	}
	return
}

// Unregister unregisters a global shortcut
func (gs *GlobalShortcuts) Unregister(accelerator string) (err error) {
	if err = gs.ctx.Err(); err != nil {
		return
	}

	// Send an event to astilectron to unregister the global shortcut
	_, err = synchronousEvent(gs.ctx, gs, gs.w, Event{Name: EventNameGlobalShortcutsCmdUnregister, TargetID: gs.id, GlobalShortcuts: &EventGlobalShortcuts{Accelerator: accelerator}}, EventNameGlobalShortcutsEventUnregistered)
	if err != nil {
		return
	}
	gs.m.Lock()
	delete(gs.callbacks, accelerator)
	gs.m.Unlock()
	return
}

// UnregisterAll unregisters all global shortcuts
func (gs *GlobalShortcuts) UnregisterAll() (err error) {
	if err = gs.ctx.Err(); err != nil {
		return
	}

	// Send an event to astilectron to unregister all global shortcuts
	_, err = synchronousEvent(gs.ctx, gs, gs.w, Event{Name: EventNameGlobalShortcutsCmdUnregisterAll, TargetID: gs.id}, EventNameGlobalShortcutsEventUnregisteredAll)
	if err != nil {
		return
	}

	gs.m.Lock()
	gs.callbacks = make(map[string]globalShortcutsCallback) // Clear the map
	gs.m.Unlock()

	return
}
