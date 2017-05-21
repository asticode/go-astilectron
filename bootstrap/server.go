package bootstrap

import (
	"net"
	"net/http"
	"text/template"

	"path/filepath"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/template"
	"github.com/julienschmidt/httprouter"
)

// Vars
var (
	templates *template.Template
)

// serve initialize an HTTP server
func serve(baseDirectoryPath string, fnR AdaptRouter, fnT TemplateData) (ln net.Listener) {
	// Init router
	var r = httprouter.New()

	// Static files
	r.ServeFiles("/static/*filepath", http.Dir(filepath.Join(baseDirectoryPath, "resources", "static")))

	// Dynamic pages
	r.GET("/templates/*page", handleTemplates(fnT))

	// Adapt router
	if fnR != nil {
		fnR(r)
	}

	// Parse templates
	var err error
	if templates, err = astitemplate.ParseDirectory(filepath.Join(baseDirectoryPath, "resources", "templates"), ".html"); err != nil {
		astilog.Fatal(err)
	}

	// Listen
	if ln, err = net.Listen("tcp", "127.0.0.1:"); err != nil {
		astilog.Fatal(err)
	}
	astilog.Debugf("Listening on %s", ln.Addr())

	// Serve
	go http.Serve(ln, r)
	return
}

// handleTemplate handles templates
func handleTemplates(fn TemplateData) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Check if template exists
		var name = p.ByName("page") + ".html"
		if templates.Lookup(name) == nil {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		// Get data
		var d interface{}
		var err error
		if fn != nil {
			if d, err = fn(name, r, p); err != nil {
				astilog.Errorf("%s while retrieving data for template %s", err, name)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		// Execute template
		if err = templates.ExecuteTemplate(rw, name, d); err != nil {
			astilog.Errorf("%s while handling template %s", err, name)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
