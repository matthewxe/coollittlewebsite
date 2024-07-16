// Setup the http server
package setup

import (
	"log"
	"net/http"

	"coollittlewebsite/internal/uno"
	"coollittlewebsite/internal/whataboutme"
)

// The port in which the server runs on
const port string = ":8080"

// Setup runs http.ListenAndServe
func Setup() {
	log.Print("Listening on " + port + "...")

	// What about me?
	whataboutme.Serve()

	// TODO: Uno 2
	uno.Serve()

	// TODO Blogs
	// blog.Serve()

	// Serve
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
