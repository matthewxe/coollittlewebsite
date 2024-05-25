package lobby

import (
	"coollittlewebsite/internal/serve/assets"
	"net/http"
)

func Serve() {

	http.HandleFunc("GET /uno", lobby)
	http.HandleFunc("GET /uno/{$}",
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/uno", http.StatusPermanentRedirect)
		})

	http.HandleFunc("GET /uno/", assets.ServeAssets)
	hub := newHub()
	go hub.run()
	http.HandleFunc("GET /uno/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	// http.HandleFunc("GET /uno/", func(w http.ResponseWriter, r *http.Request) {
	// 	assets.ServeAssets(w, r, "/uno", "/uno")
	// })
}

func lobby(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/static/uno/lobby.html")
}
