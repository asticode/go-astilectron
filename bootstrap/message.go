package bootstrap

import (
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
)

// Message represents a message
type Message struct {
	Name    string      `json:"name"`
	Payload interface{} `json:"payload"`
}

// handleMessages handles messages
func handleMessages(w *astilectron.Window, messageHandler MessageHandler) astilectron.Listener {
	return func(e astilectron.Event) (deleteListener bool) {
		// Unmarshal message
		var m Message
		var err error
		if err = e.Message.Unmarshal(&m); err != nil {
			astilog.Errorf("Unmarshaling message %+v failed", *e.Message)
			return
		}

		// Handle message
		messageHandler(w, m)
		return
	}
}
