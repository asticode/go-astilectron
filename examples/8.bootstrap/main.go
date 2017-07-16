package main

import (
	"flag"
	"os"

	"encoding/json"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/asticode/go-astilog"
)

//go:generate go-bindata -pkg $GOPACKAGE -o resources.go resources/...
func main() {
	// Parse flags
	flag.Parse()

	// Set up logger
	astilog.FlagInit()

	// Get base dir path
	var p = os.Getenv("GOPATH") + "/src/github.com/asticode/go-astilectron/examples"

	// Run bootstrap
	if err := bootstrap.Run(bootstrap.Options{
		AstilectronOptions: astilectron.Options{
			AppName:            "Astilectron",
			AppIconDefaultPath: p + "/gopher.png",
			AppIconDarwinPath:  p + "/gopher.icns",
		},
		CustomProvision: func(baseDirectoryPath string) error {
			astilog.Info("You can run your custom provisioning here!")
			return nil
		},
		Homepage: "index.html",
		MessageHandler: func(w *astilectron.Window, m bootstrap.MessageIn) {
			switch m.Name {
			case "say":
				// Unmarshal
				type P struct {
					Message string `json:"message"`
				}
				var p P
				if err := json.Unmarshal(m.Payload, &p); err != nil {
					astilog.Errorf("Unmarshaling %s failed", m.Payload)
					return
				}

				// Send
				w.Send(bootstrap.MessageOut{Name: "say", Payload: p})
			}
		},
		RestoreAssets: RestoreAssets,
		WindowOptions: &astilectron.WindowOptions{
			Center: astilectron.PtrBool(true),
			Height: astilectron.PtrInt(600),
			Width:  astilectron.PtrInt(600),
		},
	}); err != nil {
		astilog.Fatal(err)
	}
}
