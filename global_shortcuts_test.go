package astilectron

import (
	"context"
	"fmt"
	"testing"
)

func TestGlobalShortcut_Actions(t *testing.T) {
	var d = newDispatcher()
	var i = newIdentifier()
	var wrt = &mockedWriter{}
	var w = newWriter(wrt, &logger{})

	var gs = newGlobalShortcuts(context.Background(), d, i, w)

	// Register
	testObjectAction(t, func() error {
		_, e := gs.Register("Ctrl+X", func() {})
		return e
	}, gs.object, wrt, fmt.Sprintf(`{"name":"%s","targetID":"%s","globalShortcuts":{"accelerator":"Ctrl+X"}}%s`, EventNameGlobalShortcutsCmdRegister, gs.id, "\n"),
		EventNameGlobalShortcutsEventRegistered, &Event{GlobalShortcuts: &EventGlobalShortcuts{Accelerator: "Ctrl+X", IsRegistered: true}})

	// IsRegistered
	testObjectAction(t, func() error {
		_, e := gs.IsRegistered("Ctrl+Y")
		return e
	}, gs.object, wrt, fmt.Sprintf(`{"name":"%s","targetID":"%s","globalShortcuts":{"accelerator":"Ctrl+Y"}}%s`, EventNameGlobalShortcutsCmdIsRegistered, gs.id, "\n"),
		EventNameGlobalShortcutsEventIsRegistered, &Event{GlobalShortcuts: &EventGlobalShortcuts{Accelerator: "Ctrl+Y", IsRegistered: false}})

	// Unregister
	testObjectAction(t, func() error {
		return gs.Unregister("Ctrl+Z")
	}, gs.object, wrt, fmt.Sprintf(`{"name":"%s","targetID":"%s","globalShortcuts":{"accelerator":"Ctrl+Z"}}%s`, EventNameGlobalShortcutsCmdUnregister, gs.id, "\n"),
		EventNameGlobalShortcutsEventUnregistered, nil)

	// UnregisterAll
	testObjectAction(t, func() error {
		return gs.UnregisterAll()
	}, gs.object, wrt, fmt.Sprintf(`{"name":"%s","targetID":"%s"}%s`, EventNameGlobalShortcutsCmdUnregisterAll, gs.id, "\n"),
		EventNameGlobalShortcutsEventUnregisteredAll, nil)
}
