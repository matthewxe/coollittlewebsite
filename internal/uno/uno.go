package uno

import (
	"coollittlewebsite/internal/serve/assets"
	"coollittlewebsite/internal/uno/player"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"
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
	cookie, err := r.Cookie("unoName")

	if err != nil || player.PlayerList[cookie.Value].Name == "" {
		err := r.ParseForm()
		if err != nil {
			log.Fatal("Error to parse form")
			return
		}
		if r.Form.Get("name") != "" {
			key := string(randomKey(24))
			var newplyear player.Player
			newplyear.Name = r.Form.Get("name")
			player.PlayerList[key] = newplyear

			cookie := http.Cookie{}
			cookie.Name = "unoName"
			cookie.Value = key
			cookie.Expires = time.Now().Add(24 * time.Hour) // After 1 day
			// cookie.Secure = true
			cookie.Secure = false
			cookie.HttpOnly = true
			cookie.Path = "/uno"
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/uno", http.StatusSeeOther)
			return
		}
		http.ServeFile(w, r, "./web/static/uno/name.html")
		return
	}

	if r.Header.Get("getnames") == "true" {
		log.Println("yoyower")
		return
	}

	tmpl, err := template.ParseFiles("./web/static/uno/index.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	// log.Println("fuckin cookie")
	// log.Println(cookie.Value)
	// log.Println(player.PlayerList[cookie.Value].Name)

	err = tmpl.Execute(w, player.PlayerList[cookie.Value].Name)
	if err != nil {
		return
	}
	// http.ServeFile(w, r, "./web/static/uno/index.html")
}

func create(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating new lobby")
}

func randomKey(len int) (key []byte) {
	for i := 0; i < len; i++ {
		excluded := []int{1, 26, 59}
		random := randIntExclude(93, excluded)
		key = append(key, byte(random+33))
	}
	return key
}

func randIntExclude(top int, excluded []int) (random int) {
	random = (rand.Int() % top)
	for _, v := range excluded {
		if random == v {
			return randIntExclude(top, excluded)
		}
	}
	return
}
