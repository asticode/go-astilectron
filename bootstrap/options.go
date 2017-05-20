package bootstrap

import (
	"net/http"

	"github.com/asticode/go-astilectron"
	"github.com/julienschmidt/httprouter"
)

// Options represents options
type Options struct {
	AstilectronOptions astilectron.Options
	CustomProvision    CustomProvision
	Homepage           string
	RestoreAssets      RestoreAssets
	StartLoader        StartLoader
	TemplateData       TemplateData
	WindowOptions      *astilectron.WindowOptions
}

// CustomProvision is a function that executes custom provisioning
type CustomProvision func() error

// RestoreAssets is a function that restores assets namely the go-bindata's RestoreAssets method
type RestoreAssets func(dir, name string) error

// StartLoader is a function that can start a loader
type StartLoader func(a *astilectron.Astilectron)

// TemplateData is a function that retrieves a template's data
type TemplateData func(name string, r *http.Request, p httprouter.Params) (d interface{}, err error)
