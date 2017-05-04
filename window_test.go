package astilectron

import (
	"testing"

	"github.com/asticode/go-astitools/context"
	"github.com/stretchr/testify/assert"
)

func TestNewWindow(t *testing.T) {
	// Init
	a, err := New(Options{AppName: "app name", AppIconDefaultPath: "/path/to/default/icon"})
	assert.NoError(t, err)
	w, err := a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)

	// Test app name + icon
	assert.Equal(t, "app name", *w.o.Title)
	assert.Equal(t, "/path/to/default/icon", *w.o.Icon)

	// Test in display
	w, err = a.NewWindowInDisplay("http://test.com", &WindowOptions{X: PtrInt(3), Y: PtrInt(4)}, newDisplay(&DisplayOptions{Bounds: &RectangleOptions{PositionOptions: PositionOptions{X: PtrInt(1), Y: PtrInt(2)}, SizeOptions: SizeOptions{Height: PtrInt(5), Width: PtrInt(6)}}}, true))
	assert.NoError(t, err)
	assert.Equal(t, 4, *w.o.X)
	assert.Equal(t, 6, *w.o.Y)
}

func TestWindow_IsActionable(t *testing.T) {
	// Init
	a, err := New(Options{})
	assert.NoError(t, err)

	// Test canceller cancelled
	w, err := a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)
	a.canceller.Cancel()
	assert.EqualError(t, w.isActionable(), ErrCancellerCancelled.Error())

	// Test window destroyed
	w.cancel()
	assert.EqualError(t, w.isActionable(), ErrWindowDestroyed.Error())
}

func testWindowAction(t *testing.T, fn func() error, w *Window, wrt *mockedWriter, sentEvent, eventNameDone string) {
	wrt.w = []string{}
	w.c.Cancel()
	err := fn()
	assert.EqualError(t, err, ErrCancellerCancelled.Error())
	w.c = asticontext.NewCanceller()
	if eventNameDone != "" {
		wrt.fn = func() { w.d.dispatch(Event{Name: eventNameDone, TargetID: w.id}) }
	}
	err = fn()
	assert.NoError(t, err)
	assert.Equal(t, []string{sentEvent}, wrt.w)
}

func TestWindow_Actions(t *testing.T) {
	// Init
	a, err := New(Options{})
	assert.NoError(t, err)
	defer a.Close()
	go a.dispatcher.start()
	wrt := &mockedWriter{}
	a.writer = newWriter(wrt)
	w, err := a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)

	// Actions
	testWindowAction(t, func() error { return w.Blur() }, w, wrt, "{\"name\":\""+EventNameWindowCmdBlur+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventBlur)
	testWindowAction(t, func() error { return w.Center() }, w, wrt, "{\"name\":\""+EventNameWindowCmdCenter+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventMove)
	testWindowAction(t, func() error { return w.Close() }, w, wrt, "{\"name\":\""+EventNameWindowCmdClose+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventClosed)
	assert.True(t, w.isWindowDestroyed())
	w, err = a.NewWindow("http://test.com", &WindowOptions{Center: PtrBool(true)})
	assert.NoError(t, err)
	testWindowAction(t, func() error { return w.CloseDevTools() }, w, wrt, "{\"name\":\""+EventNameWindowCmdWebContentsCloseDevTools+"\",\"targetID\":\""+w.id+"\"}\n", "")
	testWindowAction(t, func() error { return w.Create() }, w, wrt, "{\"name\":\""+EventNameWindowCmdCreate+"\",\"targetID\":\""+w.id+"\",\"url\":\"http://test.com\",\"windowOptions\":{\"center\":true}}\n", EventNameWindowEventDidFinishLoad)
	testWindowAction(t, func() error { return w.Destroy() }, w, wrt, "{\"name\":\""+EventNameWindowCmdDestroy+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventClosed)
	assert.True(t, w.isWindowDestroyed())
	w, err = a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)
	testWindowAction(t, func() error { return w.Focus() }, w, wrt, "{\"name\":\""+EventNameWindowCmdFocus+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventFocus)
	testWindowAction(t, func() error { return w.Hide() }, w, wrt, "{\"name\":\""+EventNameWindowCmdHide+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventHide)
	testWindowAction(t, func() error { return w.OpenDevTools() }, w, wrt, "{\"name\":\""+EventNameWindowCmdWebContentsOpenDevTools+"\",\"targetID\":\""+w.id+"\"}\n", "")
	testWindowAction(t, func() error { return w.Maximize() }, w, wrt, "{\"name\":\""+EventNameWindowCmdMaximize+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventMaximize)
	testWindowAction(t, func() error { return w.Minimize() }, w, wrt, "{\"name\":\""+EventNameWindowCmdMinimize+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventMinimize)
	testWindowAction(t, func() error { return w.Move(3, 4) }, w, wrt, "{\"name\":\""+EventNameWindowCmdMove+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"x\":3,\"y\":4}}\n", EventNameWindowEventMove)
	var d = newDisplay(&DisplayOptions{Bounds: &RectangleOptions{PositionOptions: PositionOptions{X: PtrInt(1), Y: PtrInt(2)}, SizeOptions: SizeOptions{Height: PtrInt(1), Width: PtrInt(2)}}}, true)
	testWindowAction(t, func() error { return w.MoveInDisplay(d, 3, 4) }, w, wrt, "{\"name\":\""+EventNameWindowCmdMove+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"x\":4,\"y\":6}}\n", EventNameWindowEventMove)
	testWindowAction(t, func() error { return w.Resize(1, 2) }, w, wrt, "{\"name\":\""+EventNameWindowCmdResize+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"height\":2,\"width\":1}}\n", EventNameWindowEventResize)
	testWindowAction(t, func() error { return w.Restore() }, w, wrt, "{\"name\":\""+EventNameWindowCmdRestore+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventRestore)
	testWindowAction(t, func() error { return w.Send(true) }, w, wrt, "{\"name\":\""+EventNameWindowCmdMessage+"\",\"targetID\":\""+w.id+"\",\"message\":true}\n", "")
	testWindowAction(t, func() error { return w.Show() }, w, wrt, "{\"name\":\""+EventNameWindowCmdShow+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventShow)
	testWindowAction(t, func() error { return w.Unmaximize() }, w, wrt, "{\"name\":\""+EventNameWindowCmdUnmaximize+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventUnmaximize)
}
