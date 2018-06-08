package astilectron

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWindow(t *testing.T) {
	// Init
	a, err := New(Options{AppName: "app name", AppIconDefaultPath: "/path/to/default/icon"})
	assert.NoError(t, err)
	w, err := a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)

	// Test app name + icon
	assert.Equal(t, "app name", *w.Options.Title)
	assert.Equal(t, "/path/to/default/icon", *w.Options.Icon)

	// Test in display
	w, err = a.NewWindowInDisplay(newDisplay(&DisplayOptions{Bounds: &RectangleOptions{PositionOptions: PositionOptions{X: PtrInt(1), Y: PtrInt(2)}, SizeOptions: SizeOptions{Height: PtrInt(5), Width: PtrInt(6)}}}, true), "http://test.com", &WindowOptions{X: PtrInt(3), Y: PtrInt(4)})
	assert.NoError(t, err)
	assert.Equal(t, 4, *w.Options.X)
	assert.Equal(t, 6, *w.Options.Y)
	w, err = a.NewWindowInDisplay(newDisplay(&DisplayOptions{Bounds: &RectangleOptions{PositionOptions: PositionOptions{X: PtrInt(1), Y: PtrInt(2)}, SizeOptions: SizeOptions{Height: PtrInt(5), Width: PtrInt(6)}}}, true), "http://test.com", &WindowOptions{})
	assert.NoError(t, err)
	assert.Equal(t, 1, *w.Options.X)
	assert.Equal(t, 2, *w.Options.Y)
}

func TestWindow_Actions(t *testing.T) {
	// Init
	a, err := New(Options{})
	assert.NoError(t, err)
	defer a.Close()
	wrt := &mockedWriter{}
	a.writer = newWriter(wrt)
	w, err := a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)

	// Actions
	testObjectAction(t, func() error { return w.Blur() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdBlur+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventBlur)
	testObjectAction(t, func() error { return w.Center() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdCenter+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventMove)
	testObjectAction(t, func() error { return w.Close() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdClose+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventClosed)
	assert.True(t, w.IsDestroyed())
	w, err = a.NewWindow("http://test.com", &WindowOptions{Center: PtrBool(true)})
	assert.NoError(t, err)
	testObjectAction(t, func() error { return w.CloseDevTools() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdWebContentsCloseDevTools+"\",\"targetID\":\""+w.id+"\"}\n", "")
	testObjectAction(t, func() error { return w.Create() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdCreate+"\",\"targetID\":\""+w.id+"\",\"sessionId\":\"4\",\"url\":\"http://test.com\",\"windowOptions\":{\"center\":true}}\n", EventNameWindowEventDidFinishLoad)
	testObjectAction(t, func() error { return w.Destroy() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdDestroy+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventClosed)
	assert.True(t, w.IsDestroyed())
	w, err = a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)
	testObjectAction(t, func() error { return w.Focus() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdFocus+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventFocus)
	testObjectAction(t, func() error { return w.Hide() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdHide+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventHide)
	testObjectAction(t, func() error { return w.Log("message") }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdLog+"\",\"targetID\":\""+w.id+"\",\"message\":\"message\"}\n", "")
	testObjectAction(t, func() error { return w.Maximize() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdMaximize+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventMaximize)
	testObjectAction(t, func() error { return w.Minimize() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdMinimize+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventMinimize)
	testObjectAction(t, func() error { return w.OpenDevTools() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdWebContentsOpenDevTools+"\",\"targetID\":\""+w.id+"\"}\n", "")
	testObjectAction(t, func() error { return w.Move(3, 4) }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdMove+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"x\":3,\"y\":4}}\n", EventNameWindowEventMove)
	var d = newDisplay(&DisplayOptions{Bounds: &RectangleOptions{PositionOptions: PositionOptions{X: PtrInt(1), Y: PtrInt(2)}, SizeOptions: SizeOptions{Height: PtrInt(1), Width: PtrInt(2)}}}, true)
	testObjectAction(t, func() error { return w.MoveInDisplay(d, 3, 4) }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdMove+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"x\":4,\"y\":6}}\n", EventNameWindowEventMove)
	testObjectAction(t, func() error { return w.Resize(1, 2) }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdResize+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"height\":2,\"width\":1}}\n", EventNameWindowEventResize)
	testObjectAction(t, func() error { return w.Restore() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdRestore+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventRestore)
	testObjectAction(t, func() error { return w.Show() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdShow+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventShow)
	testObjectAction(t, func() error { return w.Unmaximize() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdUnmaximize+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventUnmaximize)
}

func TestWindow_OnLogin(t *testing.T) {
	a, err := New(Options{})
	assert.NoError(t, err)
	defer a.Close()
	wrt := &mockedWriter{wg: &sync.WaitGroup{}}
	a.writer = newWriter(wrt)
	w, err := a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)
	w.OnLogin(func(i Event) (username, password string, err error) {
		return "username", "password", nil
	})
	wrt.wg.Add(1)
	a.dispatcher.dispatch(Event{CallbackID: "1", Name: EventNameWebContentsEventLogin, TargetID: w.id})
	wrt.wg.Wait()
	assert.Equal(t, []string{"{\"name\":\"web.contents.event.login.callback\",\"targetID\":\"1\",\"callbackId\":\"1\",\"password\":\"password\",\"username\":\"username\"}\n"}, wrt.w)
}

func TestWindow_OnMessage(t *testing.T) {
	a, err := New(Options{})
	assert.NoError(t, err)
	defer a.Close()
	wrt := &mockedWriter{wg: &sync.WaitGroup{}}
	a.writer = newWriter(wrt)
	w, err := a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)
	w.OnMessage(func(m *EventMessage) interface{} {
		return "test"
	})
	wrt.wg.Add(1)
	a.dispatcher.dispatch(Event{CallbackID: "1", Name: eventNameWindowEventMessage, TargetID: w.id})
	wrt.wg.Wait()
	assert.Equal(t, []string{"{\"name\":\"window.cmd.message.callback\",\"targetID\":\"1\",\"callbackId\":\"1\",\"message\":\"test\"}\n"}, wrt.w)
}

func TestWindow_SendMessage(t *testing.T) {
	a, err := New(Options{})
	assert.NoError(t, err)
	defer a.Close()
	wrt := &mockedWriter{}
	a.writer = newWriter(wrt)
	w, err := a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)
	wrt.fn = func() {
		a.dispatcher.dispatch(Event{CallbackID: "1", Message: newEventMessage([]byte("\"bar\"")), Name: eventNameWindowEventMessageCallback, TargetID: w.id})
		wrt.fn = nil
	}
	var wg sync.WaitGroup
	wg.Add(1)
	var s string
	w.SendMessage("foo", func(m *EventMessage) {
		m.Unmarshal(&s)
		wg.Done()
	})
	wg.Wait()
	assert.Equal(t, []string{"{\"name\":\"window.cmd.message\",\"targetID\":\"1\",\"callbackId\":\"1\",\"message\":\"foo\"}\n"}, wrt.w)
	assert.Equal(t, "bar", s)
}

func TestWindow_NewMenu(t *testing.T) {
	a, err := New(Options{})
	assert.NoError(t, err)
	w, err := a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)
	m := w.NewMenu([]*MenuItemOptions{})
	assert.Equal(t, w.id, m.rootID)
}
