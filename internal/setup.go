package setup

import (
	"log"
	"net/http"

	"coollittlewebsite/internal/whataboutme"
)

const port string = ":80"

func Setup() {
	log.Print("Listening on " + port + "...")

	// What about me?
	whataboutme.Serve()

	// TODO: Blogs
	// blog.Serve()

	// TODO: Uno 2
	// uno.Serve()

	// Serve
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
