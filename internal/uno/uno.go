package uno

import (
	"coollittlewebsite/internal/serve/assets"
	"coollittlewebsite/internal/uno/lobby"
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
	http.HandleFunc("GET /uno/logout", logout)
	http.HandleFunc("GET /uno/list", list)

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
	log.Println("serving /uno")
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
			cookieNew := &http.Cookie{}
			cookieNew.Name = "unoName"
			cookieNew.Value = key
			cookieNew.Expires = time.Now().Add(365 * 24 * time.Hour) // After 1 year
			// cookie.Secure = true
			cookieNew.Secure = false
			cookieNew.HttpOnly = true
			cookieNew.Path = "/uno"
			http.SetCookie(w, cookieNew)
			cookie = cookieNew
		} else {
			http.ServeFile(w, r, "./web/static/uno/name.html")
			return
		}
	}

	tmpl, err := template.ParseFiles("./web/static/uno/index.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = tmpl.Execute(w, player.PlayerList[cookie.Value].Name)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func list(w http.ResponseWriter, r *http.Request) {
	checkForCookie(w, r)
	for i, lobbi := range lobby.LobbyList {
		idx := i + 1
		start := "<li> " + string(idx)
		w.Write([]byte(start))
		for player, _ := range lobbi.Players {
			w.Write([]byte(player.Name))
		}
		w.Write([]byte("</li>"))
	}
}

func create(w http.ResponseWriter, r *http.Request) {
	cookie := checkForCookie(w, r)
	log.Println("serving /uno/create to ", player.PlayerList[cookie.Value].Name)

	lob := lobby.NewLobby()

	lob.Leader = player.PlayerList[cookie.Value]

	log.Println(lob)
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie := checkForCookie(w, r)
	log.Println("serving /uno/logout to \"", cookie.Value, "\"")

	delete(player.PlayerList, cookie.Value)
	c := &http.Cookie{
		Name:     "unoName",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}

	http.SetCookie(w, c)
	log.Println("")
	http.Redirect(w, r, "/uno", http.StatusSeeOther)
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

func checkForCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("unoName")
	if err != nil || player.PlayerList[cookie.Value].Name == "" {
		http.Redirect(w, r, "/uno", http.StatusSeeOther)
	}
	return cookie
}
