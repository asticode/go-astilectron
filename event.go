package astilectron

import (
	"encoding/json"
	"errors"
)

// Event names
const (
	EventNameAppEventReady                     = "app.event.ready"
	EventNameAppClose                          = "app.close"
	EventNameAppCmdStop                        = "app.cmd.stop"
	EventNameAppCrash                          = "app.crash"
	EventNameAppErrorAccept                    = "app.error.accept"
	EventNameAppNoAccept                       = "app.no.accept"
	EventNameAppTooManyAccept                  = "app.too.many.accept"
	EventNameProvisionStart                    = "provision.start"
	EventNameProvisionDone                     = "provision.done"
	EventNameWindowCmdBlur                     = "window.cmd.blur"
	EventNameWindowCmdCenter                   = "window.cmd.center"
	EventNameWindowCmdClose                    = "window.cmd.close"
	EventNameWindowCmdCreate                   = "window.cmd.create"
	EventNameWindowCmdDestroy                  = "window.cmd.destroy"
	EventNameWindowCmdFocus                    = "window.cmd.focus"
	EventNameWindowCmdHide                     = "window.cmd.hide"
	EventNameWindowCmdMaximize                 = "window.cmd.maximize"
	EventNameWindowCmdMessage                  = "window.cmd.message"
	EventNameWindowCmdMinimize                 = "window.cmd.minimize"
	EventNameWindowCmdMove                     = "window.cmd.move"
	EventNameWindowCmdResize                   = "window.cmd.resize"
	EventNameWindowCmdRestore                  = "window.cmd.restore"
	EventNameWindowCmdShow                     = "window.cmd.show"
	EventNameWindowCmdUnmaximize               = "window.cmd.unmaximize"
	EventNameWindowCmdWebContentsCloseDevTools = "window.cmd.web.contents.close.dev.tools"
	EventNameWindowCmdWebContentsOpenDevTools  = "window.cmd.web.contents.open.dev.tools"
	EventNameWindowEventBlur                   = "window.event.blur"
	EventNameWindowEventClosed                 = "window.event.closed"
	EventNameWindowEventDidFinishLoad          = "window.event.did.finish.load"
	EventNameWindowEventFocus                  = "window.event.focus"
	EventNameWindowEventHide                   = "window.event.hide"
	EventNameWindowEventMaximize               = "window.event.maximize"
	EventNameWindowEventMessage                = "window.event.message"
	EventNameWindowEventMinimize               = "window.event.minimize"
	EventNameWindowEventMove                   = "window.event.move"
	EventNameWindowEventReadyToShow            = "window.event.ready.to.show"
	EventNameWindowEventResize                 = "window.event.resize"
	EventNameWindowEventRestore                = "window.event.restore"
	EventNameWindowEventShow                   = "window.event.show"
	EventNameWindowEventUnmaximize             = "window.event.unmaximize"
	EventNameWindowEventUnresponsive           = "window.event.unresponsive"
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
	Message       *EventMessage  `json:"message,omitempty"`
	URL           string         `json:"url,omitempty"`
	WindowOptions *WindowOptions `json:"windowOptions,omitempty"`
}

// EventMessage represents an event message
type EventMessage struct {
	i interface{}
}

// newEventMessage creates a new event message
func newEventMessage(i interface{}) *EventMessage {
	return &EventMessage{i: i}
}

// MarshalJSON implements the JSONMarshaler interface
func (p *EventMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.i)
}

// Unmarshal unmarshals the payload into the given interface
func (p *EventMessage) Unmarshal(i interface{}) error {
	if b, ok := p.i.([]byte); ok {
		return json.Unmarshal(b, i)
	}
	return errors.New("event message should []byte")
}

// UnmarshalJSON implements the JSONUnmarshaler interface
func (p *EventMessage) UnmarshalJSON(i []byte) error {
	p.i = i
	return nil
}
