package astilectron

import (
	"sync"
	"testing"

	"github.com/asticode/go-astikit"
	"github.com/stretchr/testify/assert"
)

func TestNewWindow(t *testing.T) {
	// Init
	a, err := New(nil, Options{AppName: "app name", AppIconDefaultPath: "/path/to/default/icon"})
	assert.NoError(t, err)
	w, err := a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)

	// Test app name + icon
	assert.Equal(t, "app name", *w.o.Title)
	assert.Equal(t, "/path/to/default/icon", *w.o.Icon)

	// Test in display
	w, err = a.NewWindowInDisplay(newDisplay(&DisplayOptions{Bounds: &RectangleOptions{PositionOptions: PositionOptions{X: astikit.IntPtr(1), Y: astikit.IntPtr(2)}, SizeOptions: SizeOptions{Height: astikit.IntPtr(5), Width: astikit.IntPtr(6)}}}, true), "http://test.com", &WindowOptions{X: astikit.IntPtr(3), Y: astikit.IntPtr(4)})
	assert.NoError(t, err)
	assert.Equal(t, 4, *w.o.X)
	assert.Equal(t, 6, *w.o.Y)
	w, err = a.NewWindowInDisplay(newDisplay(&DisplayOptions{Bounds: &RectangleOptions{PositionOptions: PositionOptions{X: astikit.IntPtr(1), Y: astikit.IntPtr(2)}, SizeOptions: SizeOptions{Height: astikit.IntPtr(5), Width: astikit.IntPtr(6)}}}, true), "http://test.com", &WindowOptions{})
	assert.NoError(t, err)
	assert.Equal(t, 1, *w.o.X)
	assert.Equal(t, 2, *w.o.Y)
}

func TestWindow_Actions(t *testing.T) {
	// Init
	a, err := New(nil, Options{})
	assert.NoError(t, err)
	defer a.Close()
	wrt := &mockedWriter{}
	a.writer = newWriter(wrt, &logger{})
	w, err := a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)
	assert.Equal(t, false, w.IsShown())

	// Actions
	testObjectAction(t, func() error { return w.Blur() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdBlur+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventBlur)
	testObjectAction(t, func() error { return w.Center() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdCenter+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventMove)
	testObjectAction(t, func() error { return w.Close() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdClose+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventClosed)
	assert.True(t, w.ctx.Err() != nil)
	w, err = a.NewWindow("http://test.com", &WindowOptions{Center: astikit.BoolPtr(true)})
	assert.NoError(t, err)
	testObjectAction(t, func() error { return w.CloseDevTools() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdWebContentsCloseDevTools+"\",\"targetID\":\""+w.id+"\"}\n", "")
	testObjectAction(t, func() error { return w.Create() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdCreate+"\",\"targetID\":\""+w.id+"\",\"sessionId\":\"4\",\"url\":\"http://test.com\",\"windowOptions\":{\"center\":true}}\n", EventNameWindowEventDidFinishLoad)
	testObjectAction(t, func() error { return w.Destroy() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdDestroy+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventClosed)
	assert.True(t, w.ctx.Err() != nil)
	w, err = a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)
	testObjectAction(t, func() error { return w.Focus() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdFocus+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventFocus)
	testObjectAction(t, func() error { return w.Hide() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdHide+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventHide)
	assert.Equal(t, false, w.IsShown())
	testObjectAction(t, func() error { return w.Log("message") }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdLog+"\",\"targetID\":\""+w.id+"\",\"message\":\"message\"}\n", "")
	testObjectAction(t, func() error { return w.Maximize() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdMaximize+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventMaximize)
	testObjectAction(t, func() error { return w.Minimize() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdMinimize+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventMinimize)
	testObjectAction(t, func() error { return w.OpenDevTools() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdWebContentsOpenDevTools+"\",\"targetID\":\""+w.id+"\"}\n", "")
	testObjectAction(t, func() error { return w.Move(3, 4) }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdMove+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"x\":3,\"y\":4}}\n", EventNameWindowEventMove)
	var d = newDisplay(&DisplayOptions{Bounds: &RectangleOptions{PositionOptions: PositionOptions{X: astikit.IntPtr(1), Y: astikit.IntPtr(2)}, SizeOptions: SizeOptions{Height: astikit.IntPtr(1), Width: astikit.IntPtr(2)}}}, true)
	testObjectAction(t, func() error { return w.MoveInDisplay(d, 3, 4) }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdMove+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"x\":4,\"y\":6}}\n", EventNameWindowEventMove)
	testObjectAction(t, func() error { return w.Resize(1, 2) }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdResize+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"height\":2,\"width\":1}}\n", EventNameWindowEventResize)
	testObjectAction(t, func() error { return w.ResizeContent(1, 2) }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdResizeContent+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"height\":2,\"width\":1}}\n", EventNameWindowEventResizeContent)
	testObjectAction(t, func() error { return w.Restore() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdRestore+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventRestore)
	testObjectAction(t, func() error { return w.SetAlwaysOnTop(true) }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdSetAlwaysOnTop+"\",\"targetID\":\""+w.id+"\",\"enable\":true}\n", EventNameWindowEventAlwaysOnTopChanged)
	testObjectAction(t, func() error { return w.Show() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdShow+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventShow)
	assert.Equal(t, true, w.IsShown())
	testObjectAction(t, func() error { return w.UpdateCustomOptions(WindowCustomOptions{HideOnClose: astikit.BoolPtr(true)}) }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdUpdateCustomOptions+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"alwaysOnTop\":true,\"height\":2,\"show\":true,\"width\":1,\"x\":4,\"y\":6,\"custom\":{\"hideOnClose\":true}}}\n", EventNameWindowEventUpdatedCustomOptions)
	testObjectAction(t, func() error { return w.Unmaximize() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdUnmaximize+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventUnmaximize)
	testObjectAction(t, func() error { return w.ExecuteJavaScript("console.log('test');") }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdWebContentsExecuteJavaScript+"\",\"targetID\":\""+w.id+"\",\"code\":\"console.log('test');\"}\n", EventNameWindowEventWebContentsExecutedJavaScript)
}

func TestWindow_OnLogin(t *testing.T) {
	a, err := New(nil, Options{})
	assert.NoError(t, err)
	defer a.Close()
	wrt := &mockedWriter{wg: &sync.WaitGroup{}}
	a.writer = newWriter(wrt, &logger{})
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
	a, err := New(nil, Options{})
	assert.NoError(t, err)
	defer a.Close()
	wrt := &mockedWriter{wg: &sync.WaitGroup{}}
	a.writer = newWriter(wrt, &logger{})
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
	a, err := New(nil, Options{})
	assert.NoError(t, err)
	defer a.Close()
	wrt := &mockedWriter{}
	a.writer = newWriter(wrt, &logger{})
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
	a, err := New(nil, Options{})
	assert.NoError(t, err)
	w, err := a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)
	m := w.NewMenu([]*MenuItemOptions{})
	assert.Equal(t, w.id, m.rootID)
}
