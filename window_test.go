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
	testObjectAction(t, func() error { return w.Blur() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdBlur+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventBlur, nil)
	testObjectAction(t, func() error { return w.Center() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdCenter+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventMove, &Event{Bounds: &RectangleOptions{
		PositionOptions: PositionOptions{X: astikit.IntPtr(3), Y: astikit.IntPtr(4)},
		SizeOptions:     SizeOptions{Height: astikit.IntPtr(1), Width: astikit.IntPtr(1)},
	}})
	bounds, err := w.Bounds()
	assert.NoError(t, err)
	assert.Equal(t, Rectangle{Position: Position{X: 3, Y: 4}, Size: Size{Height: 1, Width: 1}}, bounds)
	testObjectAction(t, func() error { return w.Close() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdClose+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventClosed, nil)
	assert.True(t, w.ctx.Err() != nil)
	w, err = a.NewWindow("http://test.com", &WindowOptions{Center: astikit.BoolPtr(true)})
	assert.NoError(t, err)
	testObjectAction(t, func() error { return w.CloseDevTools() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdWebContentsCloseDevTools+"\",\"targetID\":\""+w.id+"\"}\n", "", nil)
	testObjectAction(t, func() error { return w.Create() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdCreate+"\",\"targetID\":\""+w.id+"\",\"sessionId\":\"4\",\"url\":\"http://test.com\",\"windowOptions\":{\"center\":true}}\n", EventNameWindowEventDidFinishLoad, &Event{Bounds: &RectangleOptions{
		PositionOptions: PositionOptions{X: astikit.IntPtr(3), Y: astikit.IntPtr(4)},
		SizeOptions:     SizeOptions{Height: astikit.IntPtr(1), Width: astikit.IntPtr(1)},
	}})
	bounds, err = w.Bounds()
	assert.NoError(t, err)
	assert.Equal(t, Rectangle{Position: Position{X: 3, Y: 4}, Size: Size{Height: 1, Width: 1}}, bounds)
	testObjectAction(t, func() error { return w.Destroy() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdDestroy+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventClosed, nil)
	assert.True(t, w.ctx.Err() != nil)
	w, err = a.NewWindow("http://test.com", &WindowOptions{})
	assert.NoError(t, err)
	testObjectAction(t, func() error { return w.Focus() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdFocus+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventFocus, nil)
	testObjectAction(t, func() error { return w.Hide() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdHide+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventHide, nil)
	assert.Equal(t, false, w.IsShown())
	testObjectAction(t, func() error { return w.Log("message") }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdLog+"\",\"targetID\":\""+w.id+"\",\"message\":\"message\"}\n", "", nil)
	testObjectAction(t, func() error { return w.Maximize() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdMaximize+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventMaximize, &Event{Bounds: &RectangleOptions{
		PositionOptions: PositionOptions{X: astikit.IntPtr(0), Y: astikit.IntPtr(0)},
		SizeOptions:     SizeOptions{Height: astikit.IntPtr(100), Width: astikit.IntPtr(200)},
	}})
	bounds, err = w.Bounds()
	assert.NoError(t, err)
	assert.Equal(t, Rectangle{Position: Position{X: 0, Y: 0}, Size: Size{Height: 100, Width: 200}}, bounds)
	testObjectAction(t, func() error { return w.Minimize() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdMinimize+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventMinimize, nil)
	testObjectAction(t, func() error { return w.OpenDevTools() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdWebContentsOpenDevTools+"\",\"targetID\":\""+w.id+"\"}\n", "", nil)
	testObjectAction(t, func() error { return w.Move(3, 4) }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdMove+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"x\":3,\"y\":4}}\n", EventNameWindowEventMove, &Event{Bounds: &RectangleOptions{
		PositionOptions: PositionOptions{X: astikit.IntPtr(5), Y: astikit.IntPtr(6)},
		SizeOptions:     SizeOptions{Height: astikit.IntPtr(2), Width: astikit.IntPtr(2)},
	}})
	bounds, err = w.Bounds()
	assert.NoError(t, err)
	assert.Equal(t, Rectangle{Position: Position{X: 5, Y: 6}, Size: Size{Height: 2, Width: 2}}, bounds)
	var d = newDisplay(&DisplayOptions{Bounds: &RectangleOptions{PositionOptions: PositionOptions{X: astikit.IntPtr(1), Y: astikit.IntPtr(2)}, SizeOptions: SizeOptions{Height: astikit.IntPtr(1), Width: astikit.IntPtr(2)}}}, true)
	testObjectAction(t, func() error { return w.MoveInDisplay(d, 3, 4) }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdMove+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"x\":4,\"y\":6}}\n", EventNameWindowEventMove, &Event{Bounds: &RectangleOptions{
		PositionOptions: PositionOptions{X: astikit.IntPtr(5), Y: astikit.IntPtr(6)},
		SizeOptions:     SizeOptions{Height: astikit.IntPtr(2), Width: astikit.IntPtr(2)},
	}})
	bounds, err = w.Bounds()
	assert.NoError(t, err)
	assert.Equal(t, Rectangle{Position: Position{X: 5, Y: 6}, Size: Size{Height: 2, Width: 2}}, bounds)
	testObjectAction(t, func() error { return w.Resize(1, 2) }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdResize+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"height\":2,\"width\":1}}\n", EventNameWindowEventResize, &Event{Bounds: &RectangleOptions{
		PositionOptions: PositionOptions{X: astikit.IntPtr(5), Y: astikit.IntPtr(6)},
		SizeOptions:     SizeOptions{Height: astikit.IntPtr(4), Width: astikit.IntPtr(4)},
	}})
	bounds, err = w.Bounds()
	assert.NoError(t, err)
	assert.Equal(t, Rectangle{Position: Position{X: 5, Y: 6}, Size: Size{Height: 4, Width: 4}}, bounds)
	testObjectAction(t, func() error { return w.ResizeContent(1, 2) }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdResizeContent+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"height\":2,\"width\":1}}\n", EventNameWindowEventResizeContent, &Event{Bounds: &RectangleOptions{
		PositionOptions: PositionOptions{X: astikit.IntPtr(4), Y: astikit.IntPtr(6)},
		SizeOptions:     SizeOptions{Height: astikit.IntPtr(2), Width: astikit.IntPtr(1)},
	}})
	bounds, err = w.Bounds()
	assert.NoError(t, err)
	assert.Equal(t, Rectangle{Position: Position{X: 4, Y: 6}, Size: Size{Height: 2, Width: 1}}, bounds)
	testObjectAction(t, func() error { return w.Restore() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdRestore+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventRestore, nil)
	testObjectAction(t, func() error { return w.SetAlwaysOnTop(true) }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdSetAlwaysOnTop+"\",\"targetID\":\""+w.id+"\",\"enable\":true}\n", EventNameWindowEventAlwaysOnTopChanged, nil)
	testObjectAction(t, func() error { return w.Show() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdShow+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventShow, nil)
	assert.Equal(t, true, w.IsShown())
	testObjectAction(t, func() error { return w.UpdateCustomOptions(WindowCustomOptions{HideOnClose: astikit.BoolPtr(true)}) }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdUpdateCustomOptions+"\",\"targetID\":\""+w.id+"\",\"windowOptions\":{\"alwaysOnTop\":true,\"height\":2,\"show\":true,\"width\":1,\"x\":4,\"y\":6,\"custom\":{\"hideOnClose\":true}}}\n", EventNameWindowEventUpdatedCustomOptions, nil)
	testObjectAction(t, func() error { return w.Unmaximize() }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdUnmaximize+"\",\"targetID\":\""+w.id+"\"}\n", EventNameWindowEventUnmaximize, nil)
	testObjectAction(t, func() error { return w.ExecuteJavaScript("console.log('test');") }, w.object, wrt, "{\"name\":\""+EventNameWindowCmdWebContentsExecuteJavaScript+"\",\"targetID\":\""+w.id+"\",\"code\":\"console.log('test');\"}\n", EventNameWindowEventWebContentsExecutedJavaScript, nil)
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
