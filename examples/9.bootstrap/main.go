package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/julienschmidt/httprouter"
)

//go:generate go-bindata -pkg $GOPACKAGE -o resources.go resources/...
func main() {
	// Parse flags
	flag.Parse()

	// Set up logger
	astilog.SetLogger(astilog.New(astilog.FlagConfig()))

	// Get base dir path
	var p = os.Getenv("GOPATH") + "/src/github.com/asticode/go-astilectron/examples"

	// Run bootstrap
	if err := bootstrap.Run(bootstrap.Options{
		AdaptRouter: func(r *httprouter.Router) {
			r.GET("/custom/route", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
				w.Write([]byte("Custom route content"))
			})
		},
		AstilectronOptions: astilectron.Options{
			AppName:            "Astilectron",
			AppIconDefaultPath: p + "/gopher.png",
			AppIconDarwinPath:  p + "/gopher.icns",
		},
		CustomProvision: func(baseDirectoryPath string) error {
			astilog.Info("You can run your custom provisioning here!")
			return nil
		},
		Homepage:      "/templates/index",
		RestoreAssets: RestoreAssets,
		TemplateData: func(name string, r *http.Request, p httprouter.Params) (d interface{}, err error) {
			switch name {
			case "/index.html":
				d = struct {
					Label string
				}{Label: "Welcome to Astilectron's bootstrap!"}
			}
			return
		},
		WindowOptions: &astilectron.WindowOptions{
			Center: astilectron.PtrBool(true),
			Height: astilectron.PtrInt(600),
			Width:  astilectron.PtrInt(600),
		},
	}); err != nil {
		astilog.Fatal(err)
	}
}
