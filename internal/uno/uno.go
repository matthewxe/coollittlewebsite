package uno

import (
	"coollittlewebsite/internal/serve/assets"
	"log"
	"net/http"
)

func Serve() {
	// Main page and assets
	http.HandleFunc("GET /uno", index)
	http.HandleFunc("GET /uno/{$}",
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/uno", http.StatusPermanentRedirect)
		})
	http.HandleFunc("GET /uno/", assets.ServeAssets)

	// TODO: Creating a lobby
	http.HandleFunc("GET /uno/create", create)

	// TODO: in a lobby
	//http.HandleFunc("GET /uno/lobby", )

	// hub := newHub()
	// go hub.run()
	//
	// http.HandleFunc("GET /uno/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	serveWs(hub, w, r)
	// })
	// http.HandleFunc("GET /uno/", func(w http.ResponseWriter, r *http.Request) {
	// 	assets.ServeAssets(w, r, "/uno", "/uno")
	// })
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/static/uno/index.html")
}

func create(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating new lobby")
}
