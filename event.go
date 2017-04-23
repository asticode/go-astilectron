package astilectron

// Event names
const (
	EventNameAppEventReady           = "app.event.ready"
	EventNameAppClose                = "app.close"
	EventNameProvisionStart          = "provision.start"
	EventNameProvisionDone           = "provision.done"
	EventNameWindowCmdBlur           = "window.cmd.blur"
	EventNameWindowCmdCenter         = "window.cmd.center"
	EventNameWindowCmdClose          = "window.cmd.close"
	EventNameWindowCmdCreate         = "window.cmd.create"
	EventNameWindowCmdDestroy        = "window.cmd.destroy"
	EventNameWindowCmdFocus          = "window.cmd.focus"
	EventNameWindowCmdHide           = "window.cmd.hide"
	EventNameWindowCmdMaximize       = "window.cmd.maximize"
	EventNameWindowCmdMinimize       = "window.cmd.minimize"
	EventNameWindowCmdMove           = "window.cmd.move"
	EventNameWindowCmdResize         = "window.cmd.resize"
	EventNameWindowCmdRestore        = "window.cmd.restore"
	EventNameWindowCmdShow           = "window.cmd.show"
	EventNameWindowCmdUnmaximize     = "window.cmd.unmaximize"
	EventNameWindowDoneCreate        = "window.done.create"
	EventNameWindowEventBlur         = "window.event.blur"
	EventNameWindowEventClosed       = "window.event.closed"
	EventNameWindowEventFocus        = "window.event.focus"
	EventNameWindowEventHide         = "window.event.hide"
	EventNameWindowEventMaximize     = "window.event.maximize"
	EventNameWindowEventMinimize     = "window.event.minimize"
	EventNameWindowEventMove         = "window.event.move"
	EventNameWindowEventReadyToShow  = "window.event.ready.to.show"
	EventNameWindowEventResize       = "window.event.resize"
	EventNameWindowEventRestore      = "window.event.restore"
	EventNameWindowEventShow         = "window.event.show"
	EventNameWindowEventUnmaximize   = "window.event.unmaximize"
	EventNameWindowEventUnresponsive = "window.event.unresponsive"
)

// Other constants
const (
	mainTargetID = "main"
)

// Event represents an event
type Event struct {
	// This is the base of the event
	Name     string `json:"name"`
	TargetID string `json:"targetID"`

	// This is a list of all possible payloads.
	// A choice was made not to use interfaces since it's a pain in the ass asserting each an every payload afterwards
	// We use pointers so that omitempty works
	Message       string         `json:"message,omitempty"`
	URL           string         `json:"url,omitempty"`
	WindowOptions *WindowOptions `json:"windowOptions,omitempty"`
}
