package astilectron

// Event names
const (
	EventNameElectronReady     = "electron.ready"
	EventNameElectronStopped   = "electron.stopped"
	EventNameProvision         = "provision"
	EventNameProvisionDone     = "provision.done"
	EventNameWindowCreate      = "window.create"
	EventNameWindowCreateDone  = "window.create.done"
	EventNameWindowReadyToShow = "window.ready.to.show"
	EventNameWindowShow        = "window.show"
	EventNameWindowShowDone    = "window.show.done"
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
