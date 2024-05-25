package main

import (
	"database/sql"
	"log"
	"math/rand"
	_ "modernc.org/sqlite"
	"net/http"
	"strings"
)

const ABOUT_DIR string = "/whataboutme/"
const STATIC_DIR string = "web/static/"
const API_KEY_LENGTH = 36

var api_key = random_key()

func main() {
	db, err := sql.Open("sqlite", "./data/data.db")
	if err != nil {
		return
	}
	defer db.Close()
	setup_sql(db)

	log.Printf("The key is\n%s", api_key)
	log.Print("Listening on :8000...")

	http.HandleFunc("GET /whataboutme", func(w http.ResponseWriter, r *http.Request) { (http.ServeFile(w, r, STATIC_DIR+"hello.html")) })
	http.HandleFunc("GET "+ABOUT_DIR, serve_assets)
	http.HandleFunc("GET "+ABOUT_DIR+"blog/{id}", serve_blog)
	http.HandleFunc("GET "+ABOUT_DIR+"addanewpostyoubingus", serve_new_post)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func random_key() (key [API_KEY_LENGTH]byte) {
	for i := 0; i < API_KEY_LENGTH; i++ {
		key[i] = byte((rand.Int() % 97) + 33)
	}
	return key
}

func setup_sql(db *sql.DB) {
	table_create := `CREATE TABLE blogs (
		id INTEGER PRIMARY KEY,
		thumbnail BLOB,
		content TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)`
	db.Exec(table_create)
}

func serve_assets(w http.ResponseWriter, r *http.Request) {
	var page string = STATIC_DIR + "hello.html"
	if string(r.URL.Path) != ABOUT_DIR {
		page = STATIC_DIR + strings.TrimPrefix(r.URL.Path, ABOUT_DIR)
		log.Print("serving " + page + " hell yeah")
	}

	log.Print("serve")
	http.ServeFile(w, r, page)
}

func serve_new_post(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
func serve_blog(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
