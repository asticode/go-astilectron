package astilectron

//protocol represents the protocol object
//https://www.electronjs.org/docs/api/protocol
type Protocol struct {
	*object
}

const (
	EventNameSessionEventProtocolInterceptStringProtocol         = "protocol.event.intercept.string.protocol"
	EventNameSessionEventProtocolInterceptStringProtocolCallback = "protocol.event.intercept.string.protocol.callback"
	//todo whats the proper schema??
	EventNameSessionEventProtocolInterceptedStringProtocol = "protocol.event.intercept.string.protocol"
)

//func passed in from user
func (s *Protocol) InterceptStringProtocol(scheme string, fn func(i Event) (mimeType string, data string, deleteListener bool)) (err error) {
	// listen for this event to be called
	s.On(EventNameSessionEventProtocolInterceptStringProtocol, func(i Event) (deleteListener bool) {
		// Get mime type, data and whether the listener should be deleted.
		mimeType, data, deleteListener := fn(i)

		// Send message back
		if err = s.w.write(Event{CallbackID: i.CallbackID, Name: EventNameSessionEventProtocolInterceptStringProtocolCallback, TargetID: s.id, Scheme: scheme, MimeType: mimeType, Data: data}); err != nil {
			return
		}

		return
	})

	if err = s.ctx.Err(); err != nil {
		return
	}

	// emit this event
	_, err = synchronousEvent(s.ctx, s, s.w, Event{Name: EventNameSessionEventProtocolInterceptStringProtocol, TargetID: s.id}, EventNameSessionEventProtocolInterceptedStringProtocol)
	return
}
