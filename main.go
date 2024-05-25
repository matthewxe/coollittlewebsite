package main

import (
	"database/sql"
	"log"
	"modernc.org/sqlite"
	"net/http"
	"strings"
)

const ABOUT_DIR string = "/whataboutme/"

func main() {
	db, err := sql.Open("sqlite", "data/data.db")
	if err != nil {
		return
	}

	defer db.Close()
	sql := `CREATE TABLE analytics (
          id INTEGER PRIMARY KEY,
          user_id INT NOT NULL,  
          event_type TEXT NOT NULL,
          properties JSON NOT NULL,
          timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
       );`

	_, err := db.Exec(`CREATE TABLE blog`)
	log.Print("Listening on :8000...")

	http.HandleFunc("GET "+ABOUT_DIR, serve_about)
	log.Fatal(http.ListenAndServe(":8000", nil))

}

func serve_about(w http.ResponseWriter, r *http.Request) {
	var page string
	if string(r.URL.Path) == ABOUT_DIR {
		log.Print("serving index")
		page = "web/static/hello.html"
	} else {
		log.Print("serving " + page)
		page = "web/static/" + strings.TrimPrefix(r.URL.Path, ABOUT_DIR)
	}

	log.Print("serve")
	http.ServeFile(w, r, page)
}
