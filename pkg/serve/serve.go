// Package serve provides utilities that can be used to serve a single file
// or a directory of files to a specific path in the url
package serve

import (
	"log"
	"net/http"
	"os"
)

// StaticDir is where static files live
// Relative to the root of the repository
const StaticDir string = "web/static"

// ServeIndex serves the index.html of a given path.
// If a user enters the directory
// Write with a starting slash but no ending slash
// serve.ServeIndex("/whataboutme")
func ServeIndex(path string) {
	http.HandleFunc("GET "+path,
		func(w http.ResponseWriter, r *http.Request) {
			log.Print("Serving " + path)
			http.ServeFile(w, r, StaticDir+path+"/index.html")
		})
	// Redirect
	http.HandleFunc("GET /"+path+"/{$}",
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("redirecting %s/ -> %s", path, path)
			http.Redirect(w, r, path, http.StatusPermanentRedirect)
		})
}

// ServeAssets is meant to be used under http.HandleFunc.
// ServeAssets will try to find the file in StaticDir + r.RequestURI
// (ex. localhost/whataboutme/js/hello.js = web/static/whataboutme/js/hello.js)
// If it does not exist, then try to redirect to the root of the directory,
// If it does then serve the file.
func ServeAssets(w http.ResponseWriter, r *http.Request) {
	log.Print("serving assets " + r.RequestURI)
	var staticDir string = StaticDir + r.RequestURI

	_, err := os.Stat(staticDir)
	if os.IsNotExist(err) {
		log.Print("failed to serve asset " + staticDir)
		http.Redirect(w, r, parseRootDir(r.RequestURI), http.StatusPermanentRedirect)
		return
	}
	http.ServeFile(w, r, staticDir)
}

// Finds the root directory when given a directory
// "/go/is/the/best/language" -> "/go"
func parseRootDir(s string) (out string) {
	for i := 1; i < len(s); i++ {
		if s[i] == '/' {
			return s[:i]
		}
	}
	return s
}
