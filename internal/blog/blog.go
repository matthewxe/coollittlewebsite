package blog

import (
	"coollittlewebsite/internal/sqlite"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
)

const ABOUT_DIR string = "/whataboutme"

func ServeBlog(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
	id := r.FormValue("id")

	rows, err := sqlite.DB.Query("SELECT content from blogs WHERE id = id", id)
	if err != nil {
		return
	}
	var content string
	rows.Scan(content)

	w.Write([]byte(content))
}

func ServeNewPost(w http.ResponseWriter, r *http.Request) {
	log.Print("add a new post")
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, "web/template/newpostprompt.html")
	case http.MethodPost:
		r.ParseForm()

		log.Printf("'%s' != '%s' ? %b", string(sqlite.ApiKey), r.Form.Get("key"), string(sqlite.ApiKey) != string(r.Form.Get("key")))

		if string(sqlite.ApiKey) != string(r.Form.Get("key")) {
			log.Print("fail")
			http.Redirect(w, r, ABOUT_DIR, http.StatusFound)
			return
		}
		log.Print("fail")
		w.Write([]byte("Oogabooga"))
	}
}
