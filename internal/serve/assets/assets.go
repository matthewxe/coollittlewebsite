package assets

import (
	"log"
	"net/http"
	"os"
)

const StaticDir string = "web/static"

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

func parseRootDir(s string) (out string) {
	for i := 1; i < len(s); i++ {
		if s[i] == '/' {
			return s[:i]
		}
	}
	return s
}
