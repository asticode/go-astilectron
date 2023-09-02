package astilectron

import "context"

const (
	EventNameGlobalShortcutCmdRegister          = "global.shortcuts.event.registered"
	EventNameGlobalShortcutCmdIsRegistered      = "global.shortcuts.event.is.registered"
	EventNameGlobalShortcutCmdUnregister        = "global.shortcuts.event.unregistered"
	EventNameGlobalShortcutCmdUnregisterAll     = "global.shortcuts.event.unregistered.all"
	EventNameGlobalShortcutEventProcessFinished = "global.shortcuts.event.process.finished" // Register or Unregister process is finished
	EventNameGlobalShortcutEventTriggered       = "global.shortcuts.event.triggered"
)

type callback func()

// GlobalShortcuts represents a global shortcut
type GlobalShortcuts struct {
	*object
	callbacks map[string]*callback // Store all registered Global Shortcuts
}

func newGlobalShortcuts(ctx context.Context, d *dispatcher, i *identifier, w *writer) (gs *GlobalShortcuts) {
	var obj = newObject(ctx, d, i, w, i.new())
	gs = &GlobalShortcuts{object: obj, callbacks: make(map[string]*callback)}

	obj.On(EventNameGlobalShortcutEventTriggered, func(e Event) (deleteListener bool) { // Register the listener for the triggered event
		globalShortcutHandler(gs, e.GlobalShortcuts.Accelerator)
		return
	})
	return
}

// Register Register global shortcuts
func (gs *GlobalShortcuts) Register(accelerator string, callback callback) (isRegistered bool, err error) {

	if err = gs.ctx.Err(); err != nil {
		return
	}

	// Send an event to astilectron to register the global shortcut
	var event = Event{Name: EventNameGlobalShortcutCmdRegister, TargetID: gs.id, GlobalShortcuts: &EventGlobalShortcuts{Accelerator: accelerator}}
	result, err := synchronousEvent(gs.ctx, gs, gs.w, event, EventNameGlobalShortcutEventProcessFinished)

	if err != nil {
		return
	}

	// If registered successfully, add the callback to the map
	if result.GlobalShortcuts.IsRegistered {
		gs.callbacks[accelerator] = &callback
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
	var event = Event{Name: EventNameGlobalShortcutCmdIsRegistered, TargetID: gs.id, GlobalShortcuts: &EventGlobalShortcuts{Accelerator: accelerator}}
	result, err := synchronousEvent(gs.ctx, gs, gs.w, event, EventNameGlobalShortcutEventProcessFinished)

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
	var event = Event{Name: EventNameGlobalShortcutCmdUnregister, TargetID: gs.id, GlobalShortcuts: &EventGlobalShortcuts{Accelerator: accelerator}}
	_, err = synchronousEvent(gs.ctx, gs, gs.w, event, EventNameGlobalShortcutEventProcessFinished)

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
	var event = Event{Name: EventNameGlobalShortcutCmdUnregisterAll, TargetID: gs.id}
	_, err = synchronousEvent(gs.ctx, gs, gs.w, event, EventNameGlobalShortcutEventProcessFinished)

	if err != nil {
		return
	}

	gs.callbacks = make(map[string]*callback) // Clear the map

	return
}

// globalShortcutHandler Handle the GlobalShortcuts event triggered from astilectron
func globalShortcutHandler(gs *GlobalShortcuts, accelerator string) {
	if callback, ok := gs.callbacks[accelerator]; ok {
		(*callback)()
	}
}
