// Setup the http server
package setup

import (
	"log"
	"net/http"

	snapws "github.com/Atheer-Ganayem/SnapWS"

	"coollittlewebsite/internal/uno"
)

// "coollittlewebsite/internal/whataboutme"

// The port in which the server runs on
const (
	port        string = ":8080"
	defaultPath string = "/whataboutme"
)

var manager *snapws.RoomManager[string]

// Setup runs http.ListenAndServe
func Setup() {
	log.Print("Listening on " + port + "...")

	// What about me?
	// whataboutme.Serve()

	// Redirect to "What about me?" when the requested uri is not found
	// http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Printf("redirecting %s -> %s", r.RequestURI, defaultPath)
	// 	http.Redirect(w, r, defaultPath, http.StatusTemporaryRedirect)
	// })

	// TODO: Uno 2
	uno.Serve(&manager)
	defer manager.Shutdown()

	// TODO: Blogs
	// blog.Serve()

	// Serve
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
