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

	InitGlobalShortcuts(context.Background(), d, i, w)

	// Register
	testObjectAction(t, func() error {
		_, e := GlobalShortcutRegister("Ctrl+X", func() {})
		return e
	}, obj, wrt, fmt.Sprintf(`{"name":"%s","targetID":"%s","globalShortcut":{"accelerator":"Ctrl+X"}}%s`, EventNameGlobalShortcutCmdRegister, obj.id, "\n"),
		EventNameGlobalShortcutEventProcessFinished, &Event{GlobalShortcut: &GlobalShortcut{Accelerator: "Ctrl+X", IsRegistered: true}})

	// IsRegistered
	testObjectAction(t, func() error {
		_, e := GlobalShortcutIsRegistered("Ctrl+Y")
		return e
	}, obj, wrt, fmt.Sprintf(`{"name":"%s","targetID":"%s","globalShortcut":{"accelerator":"Ctrl+Y"}}%s`, EventNameGlobalShortcutCmdIsRegistered, obj.id, "\n"),
		EventNameGlobalShortcutEventProcessFinished, &Event{GlobalShortcut: &GlobalShortcut{Accelerator: "Ctrl+Y", IsRegistered: false}})

	// Unregister
	testObjectAction(t, func() error {
		return GlobalShortcutUnregister("Ctrl+Z")
	}, obj, wrt, fmt.Sprintf(`{"name":"%s","targetID":"%s","globalShortcut":{"accelerator":"Ctrl+Z"}}%s`, EventNameGlobalShortcutCmdUnregister, obj.id, "\n"),
		EventNameGlobalShortcutEventProcessFinished, nil)

	// UnregisterAll
	testObjectAction(t, func() error {
		return GlobalShortcutUnregisterAll()
	}, obj, wrt, fmt.Sprintf(`{"name":"%s","targetID":"%s"}%s`, EventNameGlobalShortcutCmdUnregisterAll, obj.id, "\n"),
		EventNameGlobalShortcutEventProcessFinished, nil)
}
