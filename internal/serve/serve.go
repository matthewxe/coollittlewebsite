package serve

import (
	"log"
	"net/http"

	"coollittlewebsite/internal/serve/assets"
	"coollittlewebsite/internal/uno"
	"coollittlewebsite/internal/webhooks"
)

const port string = ":80"

func Setup() {
	log.Print("Listening on " + port + "...")

	// Handle /whataboutme
	http.HandleFunc("GET /whataboutme",
		func(w http.ResponseWriter, r *http.Request) {
			log.Print("Serving /whataboutme")
			http.ServeFile(w, r, "web/static/whataboutme/index.html")
		})
	http.HandleFunc("GET /whataboutme/{$}",
		func(w http.ResponseWriter, r *http.Request) {
			log.Print("redirecting /whataboutme/ -> /whataboutme")
			http.Redirect(w, r, "/whataboutme", http.StatusPermanentRedirect)
		})
	// Assets for /whataboutme
	http.HandleFunc("GET /whataboutme/", assets.ServeAssets)

	// Github webhook primitive CI/CD
	http.HandleFunc("POST "+webhooks.GithubWebhookDir, webhooks.GithubWebhookHTTP)

	// TODO: Blogs
	// http.HandleFunc("GET "+ABOUT_DIR+"/blog/{id}", blog.ServeBlog)
	// http.HandleFunc("GET "+ABOUT_DIR+"/addanewpostyoubingus", blog.ServeNewPost)
	// http.HandleFunc("POST "+ABOUT_DIR+"/addanewpostyoubingus", blog.ServeNewPost)

	// TODO: Uno 2
	uno.Serve()

	// Serve
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
