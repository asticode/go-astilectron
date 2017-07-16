package astilectron

import (
	"testing"

	"github.com/asticode/go-astitools/context"
	"github.com/stretchr/testify/assert"
)

func TestTray_Actions(t *testing.T) {
	// Init
	var c = asticontext.NewCanceller()
	var d = newDispatcher()
	go d.start()
	var i = newIdentifier()
	var wrt = &mockedWriter{}
	var w = newWriter(wrt)
	var tr = newTray(&TrayOptions{
		Image:   PtrStr("/path/to/image"),
		Tooltip: PtrStr("tooltip"),
	}, c, d, i, w)

	// Actions
	testObjectAction(t, func() error { return tr.Create() }, tr.object, wrt, "{\"name\":\""+EventNameTrayCmdCreate+"\",\"targetID\":\""+tr.id+"\",\"trayOptions\":{\"image\":\"/path/to/image\",\"tooltip\":\"tooltip\"}}\n", EventNameTrayEventCreated)
	testObjectAction(t, func() error { return tr.Destroy() }, tr.object, wrt, "{\"name\":\""+EventNameTrayCmdDestroy+"\",\"targetID\":\""+tr.id+"\"}\n", EventNameTrayEventDestroyed)
	assert.True(t, tr.IsDestroyed())
}
