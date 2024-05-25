package main

import (
	"database/sql"
	"log"
	"math/rand"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"strings"
	// "text/template"
)

const ABOUT_DIR string = "/whataboutme"
const STATIC_DIR string = "web/static"

var api_key = random_key(32)
var db *sql.DB

func main() {
	db, err := sql.Open("sqlite", "./data/data.db")
	if err != nil {
		return
	}
	defer db.Close()
	setup_sql(db)

	log.Printf("The key is\n%s", api_key)
	log.Print("Listening on :8000...")

	http.HandleFunc("GET "+ABOUT_DIR,
		func(w http.ResponseWriter, r *http.Request) {
			log.Print("serving about me")
			http.ServeFile(w, r, STATIC_DIR+"/index.html")
		})
	http.HandleFunc("GET "+ABOUT_DIR+"/", serve_assets)
	http.HandleFunc("GET "+ABOUT_DIR+"/blog/{id}", serve_blog)
	http.HandleFunc("GET "+ABOUT_DIR+"/addanewpostyoubingus", serve_new_post)
	http.HandleFunc("POST "+ABOUT_DIR+"/addanewpostyoubingus", serve_new_post)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func random_key(len int) (key []byte) {
	for i := 0; i < len; i++ {
		key = append(key, byte((rand.Int()%97)+33))
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
	var page string = strings.TrimPrefix(r.RequestURI, ABOUT_DIR)
	if page == "/" {
		http.Redirect(w, r, ABOUT_DIR, http.StatusPermanentRedirect)
	}
	page = STATIC_DIR + page
	log.Print("serving asset " + page)

	_, err := os.Stat(page)
	if os.IsNotExist(err) {
		log.Print("failed to serve asset " + page)
		http.Redirect(w, r, ABOUT_DIR, http.StatusPermanentRedirect)
		return
	}
	http.ServeFile(w, r, page)
}

func serve_new_post(w http.ResponseWriter, r *http.Request) {
	log.Print("add a new post")
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, "web/template/newpostprompt.html")
	case http.MethodPost:
		r.ParseForm()

		log.Printf("'%s' != '%s' ? %b", string(api_key), r.Form.Get("key"), string(api_key) != string(r.Form.Get("key")))

		if string(api_key) != string(r.Form.Get("key")) {
			log.Print("fail")
			http.Redirect(w, r, ABOUT_DIR, http.StatusFound)
			return
		}
			log.Print("fail")
		w.Write([]byte("Oogabooga"))
	}
}
func serve_blog(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
	id := r.FormValue("id")

	rows, err := db.Query("SELECT content from blogs WHERE id = id", id)
	if err != nil {
		return
	}
	var content string
	rows.Scan(content)

	w.Write([]byte(content))
}
