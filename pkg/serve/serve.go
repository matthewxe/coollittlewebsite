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

// RedirectSlash sets up a http.HandleFunc that will automatically redirect
// users to a non slash ending path of a path
//
// this/is/a/path/ -> this/is/a/path
func RedirectSlash(path string) {
	http.HandleFunc("GET "+path+"/{$}",
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("redirecting %s/ -> %s", path, path)
			http.Redirect(w, r, path, http.StatusPermanentRedirect)
		})
}

// ServeFile takes a path string and dir string. Path is where a file should
// Appear on the website path. file is where a file should be looked for relative
// to StaticDir.
// Write with a starting slash but no ending slash
//
// ex. serve.ServeFile("/whataboutme" "/whataboutme/favicon.icon")
func ServeFile(path string, file string) {
	http.HandleFunc("GET "+path,
		func(w http.ResponseWriter, r *http.Request) {
			log.Print("Serving " + path)
			http.ServeFile(w, r, StaticDir+file)
		})
	// Redirect
	RedirectSlash(path)
}

// ServeIndex serves the index.html of a given path.
// It is simply a remap of ServeFile.
//
// ex. serve.ServeIndex("/whataboutme")
func ServeIndex(path string) {
	ServeFile(path, path+"/index.html")
}

// ServeAssets is meant to be used under http.HandleFunc.
// ServeAssets will try to find the file in StaticDir + r.RequestURI.
// If it does not exist, then try to redirect to the root of the directory.
// If it does then serve the file.
//
// ex. localhost/whataboutme/js/hello.js -> web/static/whataboutme/js/hello.js
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
//
// "/go/is/the/best/language" -> "/go"
func parseRootDir(s string) (out string) {
	for i := 1; i < len(s); i++ {
		if s[i] == '/' {
			return s[:i]
		}
	}
	return s
}
