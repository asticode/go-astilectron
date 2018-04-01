package astilectron

import (
	"testing"

	"github.com/asticode/go-astitools/context"
)

func TestNotification_Actions(t *testing.T) {
	// Init
	var c = asticontext.NewCanceller()
	var d = newDispatcher()
	var i = newIdentifier()
	var wrt = &mockedWriter{}
	var w = newWriter(wrt)
	var n = newNotification(&NotificationOptions{
		Body:             "body",
		HasReply:         PtrBool(true),
		Icon:             "/path/to/icon",
		ReplyPlaceholder: "placeholder",
		Silent:           PtrBool(true),
		Sound:            "sound",
		Subtitle:         "subtitle",
		Title:            "title",
	}, true, c, d, i, w)

	// Actions
	testObjectAction(t, func() error { return n.Create() }, n.object, wrt, "{\"name\":\""+eventNameNotificationCmdCreate+"\",\"targetID\":\""+n.id+"\",\"notificationOptions\":{\"body\":\"body\",\"hasReply\":true,\"icon\":\"/path/to/icon\",\"replyPlaceholder\":\"placeholder\",\"silent\":true,\"sound\":\"sound\",\"subtitle\":\"subtitle\",\"title\":\"title\"}}\n", EventNameNotificationEventCreated)
	testObjectAction(t, func() error { return n.Show() }, n.object, wrt, "{\"name\":\""+eventNameNotificationCmdShow+"\",\"targetID\":\""+n.id+"\"}\n", EventNameNotificationEventShown)
}
