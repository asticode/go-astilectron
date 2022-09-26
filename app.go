package astilectron

const (
	appCmdSetAsDefaultProtocolClient      = "app.cmd.setas.default.protocol.client"
	appEventSetAsDefaultProtocolClient    = "app.event.setas.default.protocol.client"
	appCmdRemoveAsDefaultProtocolClient   = "app.cmd.removeas.default.protocol.client"
	appEventRemoveAsDefaultProtocolClient = "app.event.removeas.default.protocol.client"
	appCmdIsDefaultProtocolClient         = "app.cmd.is.default.protocol.client"
	appEventIsDefaultProtocolClient       = "app.event.is.default.protocol.client"
)

// Set app as default protocol handler
func (a *Astilectron) SetAsDefaultProtocolClient(protocol string, path string, args []string) (success bool, err error) {
	var e Event
	success = false

	event := Event{Name: appCmdSetAsDefaultProtocolClient,
		Protocol: protocol}
	if path != "" {
		event.Path = path
	}
	if len(args) > 0 {
		event.Args = args
	}

	if e, err = synchronousEvent(a.worker.Context(), a, a.writer,
		event, appEventSetAsDefaultProtocolClient); err != nil {
		return
	}

	if e.Success != nil {
		success = *e.Success
	}
	return
}

// Remove app as default protocol handler
func (a *Astilectron) RemoveAsDefaultProtocolClient(protocol string, path string, args []string) (success bool, err error) {
	var e Event
	success = false

	event := Event{Name: appCmdRemoveAsDefaultProtocolClient,
		Protocol: protocol}
	if path != "" {
		event.Path = path
	}
	if len(args) > 0 {
		event.Args = args
	}

	if e, err = synchronousEvent(a.worker.Context(), a, a.writer,
		event,
		appEventRemoveAsDefaultProtocolClient); err != nil {
		return
	}

	if e.Success != nil {
		success = *e.Success
	}
	return
}

// Check if app is the default protocol handler
func (a *Astilectron) IsDefaultProtocolClient(protocol string, path string, args []string) (success bool, err error) {
	var e Event
	success = false

	event := Event{Name: appCmdIsDefaultProtocolClient,
		Protocol: protocol}
	if path != "" {
		event.Path = path
	}
	if len(args) > 0 {
		event.Args = args
	}

	if e, err = synchronousEvent(a.worker.Context(), a, a.writer,
		event, appEventIsDefaultProtocolClient); err != nil {
		return
	}

	if e.Success != nil {
		success = *e.Success
	}
	return
}
