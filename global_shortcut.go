package astilectron

import (
	"context"
	"errors"
	"strings"
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
	Accelerator  string    `json:"accelerator,omitempty"`  // Accelerator of the global globalShortcuts
	IsRegistered bool      `json:"isRegistered,omitempty"` // Whether the global shortcut is registered
	Callback     *callback `json:"-"`
}
type callback func()

var gss = make(map[string]*GlobalShortcut) // Store all registered Global Shortcuts
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
func GlobalShortcutRegister(accelerator string, callback callback) bool {

	// Check if the accelerator is valid
	keySet, err := parseAccelerator(accelerator)
	if err != nil || len(keySet) == 0 {
		obj.w.l.Error(err)
		return false
	}

	// If CmdOrCtrl is used, replace it with Cmd and Ctrl
	if _, ok := keySet["commandorcontrol"]; ok {

		delete(keySet, "commandorcontrol")

		// Replace CmdOrCtrl with Cmd
		keySet["command"] = true
		var result1 = GlobalShortcutRegister(keySetToString(keySet), callback)

		// Replace CmdOrCtrl with Ctrl
		delete(keySet, "command")
		keySet["control"] = true
		var result2 = GlobalShortcutRegister(keySetToString(keySet), callback)

		return result1 || result2 // Either Cmd or Ctrl is registered successfully
	}

	var gs = GlobalShortcut{Accelerator: accelerator, Callback: &callback, object: obj}

	// Send an event to astilectron to register the global shortcut
	var event = Event{Name: EventNameGlobalShortcutCmdRegister, TargetID: gs.id, GlobalShortcut: &gs}
	result, err := synchronousEvent(gs.ctx, gs, gs.w, event, EventNameGlobalShortcutEventProcessFinished)

	if err != nil {
		return false
	}

	// If registered successfully, add this `GlobalShortcut` object to the map
	if result.GlobalShortcut.IsRegistered {
		gss[accelerator] = &gs
	}

	return result.GlobalShortcut.IsRegistered
}

// GlobalShortcutIsRegistered Check if a global shortcut is registered
func GlobalShortcutIsRegistered(accelerator string) bool {
	var gs = GlobalShortcut{Accelerator: accelerator, object: obj}

	// Send an event to astilectron to check if global shortcut is registered
	var event = Event{Name: EventNameGlobalShortcutCmdIsRegistered, TargetID: gs.id, GlobalShortcut: &gs}
	result, err := synchronousEvent(gs.ctx, gs, gs.w, event, EventNameGlobalShortcutEventProcessFinished)

	if err != nil {
		return false
	}

	return result.GlobalShortcut.IsRegistered
}

// GlobalShortcutUnregister Unregister a global shortcut
func GlobalShortcutUnregister(accelerator string) {

	// Remove the callback from the map and remove the listener
	keySetToUnregister, err := parseAccelerator(accelerator) // Accelerator to unregister
	if err != nil || len(keySetToUnregister) == 0 {
		obj.w.l.Error(err)
		return
	}

	// If CmdOrCtrl is used, replace it with Cmd and Ctrl
	if _, ok := keySetToUnregister["commandorcontrol"]; ok {

		delete(keySetToUnregister, "commandorcontrol")

		// Replace CmdOrCtrl with Cmd
		keySetToUnregister["command"] = true
		GlobalShortcutUnregister(keySetToString(keySetToUnregister))

		// Replace CmdOrCtrl with Ctrl
		delete(keySetToUnregister, "command")
		keySetToUnregister["control"] = true
		GlobalShortcutUnregister(keySetToString(keySetToUnregister))

		return
	}

	// Iterate through all registered global shortcuts to find the equivalent one
	for acc, gs := range gss {
		keySetRegistered, _ := parseAccelerator(acc)

		var isEqual = isKeySetsEqual(keySetToUnregister, keySetRegistered)
		if isEqual { // Found the equivalent global shortcut

			// Send an event to astilectron to unregister the global shortcut
			var event = Event{Name: EventNameGlobalShortcutCmdUnregister, TargetID: gs.id, GlobalShortcut: gs}
			_, err := synchronousEvent(gs.ctx, gs, gs.w, event, EventNameGlobalShortcutEventProcessFinished)

			if err != nil {
				return
			}
			delete(gss, acc) // Remove this `GlobalShortcut` from the map

			break
		}
	}
}

// GlobalShortcutUnregisterAll Unregister all global shortcuts
func GlobalShortcutUnregisterAll() {

	// Send an event to astilectron to unregister all global shortcuts
	var event = Event{Name: EventNameGlobalShortcutCmdUnregisterAll, TargetID: obj.id}
	_, err := synchronousEvent(obj.ctx, obj, obj.w, event, EventNameGlobalShortcutEventProcessFinished)

	if err != nil {
		return
	}

	gss = make(map[string]*GlobalShortcut) // Clear the map
}

// parseAccelerator Parse the accelerator string into a key set
// Note that the boolean value of returned map is meaningless. We just care if the key exists
func parseAccelerator(accelerator string) (map[string]bool, error) {
	var tokens = []string{
		"Command", "Control", "CommandOrControl", "Alt",
		"Option", "AltGr", "Shift", "Super", "Meta",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
		"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T",
		"U", "V", "W", "X", "Y", "Z",
		"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8",
		"F9", "F10", "F11", "F12", "F13", "F14", "F15",
		"F16", "F17", "F18", "F19", "F20", "F21", "F22", "F23", "F24",
		")", "!", "@", "#", "$", "%", "^", "&", "*", "(",
		":", ";", ":", "+", "=", "<", ",", "_", "-", ">",
		".", "?", "/", "~", "`", "{", "]", "[", "|", "\\",
		"}", "\"",
		"Plus", "Space", "Tab", "Capslock", "Numlock",
		"Scrolllock", "Backspace", "Delete", "Insert",
		"Enter", "Up", "Down", "Left", "Right",
		"Home", "End", "PageUp", "PageDown", "Escape",
		"VolumeUp", "VolumeDown", "VolumeMute",
		"MediaNextTrack", "MediaPreviousTrack", "MediaStop",
		"MediaPlayPause", "PrintScreen",
		"Num0", "Num1", "Num2", "Num3", "Num4",
		"Num5", "Num6", "Num7", "Num8", "Num9",
		"NumDec", "NumAdd", "NumSub",
		"NumMul", "NumDiv",
	}

	var tokenAlias = map[string]string{
		"cmd":       "Command",
		"ctrl":      "Control",
		"cmdorctrl": "CommandOrControl",
		"return":    "Enter",
		"esc":       "Escape",
	}

	accelerator = strings.TrimSpace(accelerator)
	var keys = strings.Split(accelerator, "+")
	keySet := make(map[string]bool)

	for _, key := range keys {

		key = strings.ToLower(key)

		if alias, ok := tokenAlias[key]; ok { // The key is an alias of token
			keySet[strings.ToLower(alias)] = true

		} else if isInSlice(tokens, key) { // The key is a token
			keySet[key] = true

		} else {
			return nil, errors.New("invalid accelerator: " + accelerator)
		}
	}

	// If CmdOrCtrl and Cmd are both in the key set, remove CmdOrCtrl
	// If CmdOrCtrl and Ctrl are both in the key set, remove CmdOrCtrl
	if keySet["commandorcontrol"] && keySet["command"] {
		delete(keySet, "commandorcontrol")
	} else if keySet["commandorcontrol"] && keySet["control"] {
		delete(keySet, "commandorcontrol")
	}

	return keySet, nil
}

// Key set to accelerator string
func keySetToString(keySet map[string]bool) string {
	var keys []string
	for k := range keySet {
		keys = append(keys, k)
	}
	return strings.Join(keys, "+")
}

// Check if two key sets are equivalent
// Assert that the CmdOrCtrl key doesn't exist in both key sets
func isKeySetsEqual(keySetA map[string]bool, keySetB map[string]bool) bool {

	if len(keySetA) != len(keySetB) {
		return false
	}

	for k, _ := range keySetA {
		if _, ok := keySetB[k]; !ok {
			return false
		}
	}
	return true
}

// isInSlice Checks if an item is in a slice
func isInSlice(slice []string, item string) bool {
	for _, s := range slice {
		if strings.ToLower(s) == strings.ToLower(item) {
			return true
		}
	}
	return false
}

// globalShortcutHandler Handle the GlobalShortcut event triggered from astilectron
func globalShortcutHandler(accelerator string) {
	if gs, ok := gss[accelerator]; ok {
		(*(*gs).Callback)()
	}
}
