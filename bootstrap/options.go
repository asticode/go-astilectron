package bootstrap

import (
	"net/http"

	"github.com/asticode/go-astilectron"
	"github.com/julienschmidt/httprouter"
)

// Options represents options
type Options struct {
	AdaptAstilectron   AdaptAstilectron
	AdaptRouter        AdaptRouter
	AdaptWindow        AdaptWindow
	AstilectronOptions astilectron.Options
	BaseDirectoryPath  string
	CustomProvision    CustomProvision
	Debug              bool
	Homepage           string
	MessageHandler     MessageHandler
	RestoreAssets      RestoreAssets
	TemplateData       TemplateData
	WindowOptions      *astilectron.WindowOptions
}

// AdaptAstilectron is a function that adapts astilectron
type AdaptAstilectron func(a *astilectron.Astilectron)

// AdaptRouter is a function that adapts the router
type AdaptRouter func(r *httprouter.Router)

// AdaptWindow is a function that adapts the window
type AdaptWindow func(w *astilectron.Window)

// CustomProvision is a function that executes custom provisioning
type CustomProvision func(baseDirectoryPath string) error

// MessageHandler is a functions that handles messages
type MessageHandler func(w *astilectron.Window, m MessageIn)

// RestoreAssets is a function that restores assets namely the go-bindata's RestoreAssets method
type RestoreAssets func(dir, name string) error

// TemplateData is a function that retrieves a template's data
type TemplateData func(name string, r *http.Request, p httprouter.Params) (d interface{}, err error)
