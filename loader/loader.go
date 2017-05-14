package astiloader

import (
	"sync"

	"github.com/andlabs/ui"
	"github.com/asticode/go-astilectron"
)

// Loader represents a loader
type Loader struct {
	d int
	m *sync.Mutex
	p *ui.ProgressBar
	t int
	w *ui.Window
}

// New creates a new loader
func New() *Loader {
	return &Loader{m: &sync.Mutex{}}
}

// NewForAstilectron creates a new loader for Astilectron
func NewForAstilectron(a *astilectron.Astilectron) (l *Loader) {
	l = New().Add(7)
	a.On(astilectron.EventNameProvisionAstilectronAlreadyProvisioned, func(e astilectron.Event) (deleteListener bool) {
		l.Done(3)
		return
	})
	a.On(astilectron.EventNameProvisionAstilectronFinished, func(e astilectron.Event) (deleteListener bool) {
		l.Done(1)
		return
	})
	a.On(astilectron.EventNameProvisionAstilectronMoved, func(e astilectron.Event) (deleteListener bool) {
		l.Done(1)
		return
	})
	a.On(astilectron.EventNameProvisionAstilectronUnzipped, func(e astilectron.Event) (deleteListener bool) {
		l.Done(1)
		return
	})
	a.On(astilectron.EventNameProvisionElectronAlreadyProvisioned, func(e astilectron.Event) (deleteListener bool) {
		l.Done(3)
		return
	})
	a.On(astilectron.EventNameProvisionElectronFinished, func(e astilectron.Event) (deleteListener bool) {
		l.Done(1)
		return
	})
	a.On(astilectron.EventNameProvisionElectronMoved, func(e astilectron.Event) (deleteListener bool) {
		l.Done(1)
		return
	})
	a.On(astilectron.EventNameProvisionElectronUnzipped, func(e astilectron.Event) (deleteListener bool) {
		l.Done(1)
		return
	})
	a.On(astilectron.EventNameAppEventReady, func(e astilectron.Event) (deleteListener bool) {
		l.Done(1)
		return
	})
	return
}

// Add adds n new steps
func (l *Loader) Add(n int) *Loader {
	l.m.Lock()
	defer l.m.Unlock()
	l.t += n
	return l
}

// Done signifies n steps are done
func (l *Loader) Done(n int) {
	l.m.Lock()
	defer l.m.Unlock()
	l.d += n
	if l.d >= l.t {
		l.Stop()
	} else if l.p != nil {
		ui.QueueMain(func() { l.p.SetValue(int(100 * l.d / l.t)) })
	}
}

// Start starts the loader
func (l *Loader) Start() error {
	return ui.Main(func() {
		l.p = ui.NewProgressBar()
		var box = ui.NewVerticalBox()
		box.SetPadded(true)
		box.Append(l.p, true)
		l.w = ui.NewWindow("Loading...", 300, 1, false)
		l.w.SetMargined(true)
		l.w.SetChild(box)
		l.w.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		l.w.Show()
	})
}

// stop stops the loader
func (l *Loader) Stop() {
	if l.w != nil {
		ui.QueueMain(func() {
			l.w.Destroy()
			l.w = nil
		})
	}
}
