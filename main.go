package main

import (
	"log"
	"net/http"
	"strings"
)

const DIR string = "/whatisaboutme/"

func main() {

	log.Print("Listening on :8000...")
	http.HandleFunc(DIR+"*", serve_page)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func serve_page(w http.ResponseWriter, r *http.Request) {
	var page string
	if string(r.URL.Path) == DIR {
		log.Print("serving index")
		page = "web/static/hello.html"
	} else if strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	} else {
		page = "web/static/" + strings.TrimPrefix(r.URL.Path, DIR)
		log.Print("serving " + page)
	}

	http.ServeFile(w, r, page)
}
