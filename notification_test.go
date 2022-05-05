package astilectron

import (
	"context"
	"testing"

	"github.com/asticode/go-astikit"
)

func TestNotification_Actions(t *testing.T) {
	// Init
	var d = newDispatcher()
	var i = newIdentifier()
	var wrt = &mockedWriter{}
	var w = newWriter(wrt, &logger{})
	var n = newNotification(context.Background(), &NotificationOptions{
		Body:             "body",
		HasReply:         astikit.BoolPtr(true),
		Icon:             "/path/to/icon",
		ReplyPlaceholder: "placeholder",
		Silent:           astikit.BoolPtr(true),
		Sound:            "sound",
		Subtitle:         "subtitle",
		Title:            "title",
	}, true, d, i, w)

	// Actions
	testObjectAction(t, func() error { return n.Create() }, n.object, wrt, "{\"name\":\""+eventNameNotificationCmdCreate+"\",\"targetID\":\""+n.id+"\",\"notificationOptions\":{\"body\":\"body\",\"hasReply\":true,\"icon\":\"/path/to/icon\",\"replyPlaceholder\":\"placeholder\",\"silent\":true,\"sound\":\"sound\",\"subtitle\":\"subtitle\",\"title\":\"title\"}}\n", EventNameNotificationEventCreated, nil)
	testObjectAction(t, func() error { return n.Show() }, n.object, wrt, "{\"name\":\""+eventNameNotificationCmdShow+"\",\"targetID\":\""+n.id+"\"}\n", EventNameNotificationEventShown, nil)
}
