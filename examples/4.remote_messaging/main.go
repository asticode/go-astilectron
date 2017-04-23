package main

import (
	"flag"
	"net/http"
	"os"

	"time"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

func main() {
	// Parse flags
	flag.Parse()

	// Set up logger
	astilog.SetLogger(astilog.New(astilog.FlagConfig()))

	// Start server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<!DOCTYPE html>
		<html lang="en">
		<head>
		    <meta charset="UTF-8">
		    <title>Hello world</title>
		</head>
		<body>
		    <span id="message">Hello world</span>
		    <script>
		    	document.addEventListener('astilectron-ready', function() {
		    	    astilectron.listen(function(message) {
			        document.getElementById('message').innerHTML = message
			        astilectron.send("I'm good bro")
			    });
		    	})
		    </script>
		</body>
		</html>`))
	})
	go http.ListenAndServe("127.0.0.1:4000", nil)

	// Create astilectron
	var a *astilectron.Astilectron
	var err error
	if a, err = astilectron.New(astilectron.Options{BaseDirectoryPath: os.Getenv("GOPATH") + "/src/github.com/asticode/go-astilectron/examples"}); err != nil {
		astilog.Fatal(errors.Wrap(err, "creating new astilectron failed"))
	}
	defer a.Close()
	a.HandleSignals()
	a.On(astilectron.EventNameAppClose, func(e astilectron.Event) (deleteListener bool) {
		a.Stop()
		return
	})

	// Start
	if err = a.Start(); err != nil {
		astilog.Fatal(errors.Wrap(err, "starting failed"))
	}

	// Create window
	var w *astilectron.Window
	if w, err = a.NewWindow("http://127.0.0.1:4000", &astilectron.WindowOptions{
		Center: astilectron.PtrBool(true),
		Height: astilectron.PtrInt(600),
		Width:  astilectron.PtrInt(600),
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "new window failed"))
	}
	w.On(astilectron.EventNameWindowEventMessage, func(e astilectron.Event) (deleteListener bool) {
		var m string
		e.Message.Unmarshal(&m)
		astilog.Infof("Received message %s", m)
		return
	})
	w.Create()

	// Send message
	time.Sleep(time.Second)
	w.Send("What's up?")

	// Blocking pattern
	a.Wait()
}
