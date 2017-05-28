package bootstrap

import (
	"encoding/json"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
)

// MessageOut represents a message going out
type MessageOut struct {
	Name    string      `json:"name"`
	Payload interface{} `json:"payload"`
}

// MessageIn represents a message going in
type MessageIn struct {
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload"`
}

// handleMessages handles messages
func handleMessages(w *astilectron.Window, messageHandler MessageHandler) astilectron.Listener {
	return func(e astilectron.Event) (deleteListener bool) {
		// Unmarshal message
		var m MessageIn
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
