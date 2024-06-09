package uno

import (
	"coollittlewebsite/internal/serve/assets"
	"coollittlewebsite/internal/uno/lobby"
	"coollittlewebsite/internal/uno/player"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func Serve() {
	// Main page and assets
	http.HandleFunc("GET /uno", serveIndex)
	http.HandleFunc("GET /uno/{$}",
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/uno", http.StatusPermanentRedirect)
		})
	http.HandleFunc("GET /uno/", assets.ServeAssets)

	// Logging out
	http.HandleFunc("GET /uno/logout", serveLogout)

	// TODO: Creating a lobby
	http.HandleFunc("GET /uno/create", serveCreate)
	http.HandleFunc("GET /uno/list", serveList)

	// TODO: Serve a lobby
	http.HandleFunc("GET /uno/lobby/{id}", serveLobby)
	http.HandleFunc("GET /uno/lobby/{id}/ws", serveLobby)

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

func serveIndex(w http.ResponseWriter, r *http.Request) { // {{{
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
} // }}}

func serveList(w http.ResponseWriter, r *http.Request) { // {{{
	checkForCookie(w, r)
	var ready string
	var ongoing string
	var done string
	for i, lobbi := range lobby.LobbyList {
		out := fmt.Sprintf("<li>%v. ", i+1)
		out += lobbi.Leader.Name + "(leader)"
		for player := range lobbi.Players {
			out += ", "
			out += player.Name
		}
		out += fmt.Sprintf("  <button onmousedown=\"window.location.href = '/uno/lobby/ %v';\">Join</button></li>", i)

		switch lobbi.State {
		case 0:
			ready += out
		case 1:
			ongoing += out
		case 2:
			done += out
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write([]byte("<h1>Ready</h1>" + ready + "<h1>Ongoing</h1>" + ongoing + "<h1>Done</h1>" + done))

	if err != nil {
		log.Fatal(err)
		return
	}
} // }}}

func serveCreate(w http.ResponseWriter, r *http.Request) { // {{{
	cookie := checkForCookie(w, r)
	log.Println("serving /uno/create to ", player.PlayerList[cookie.Value].Name)

	lobbyId := strconv.Itoa(lobby.NewLobby(player.PlayerList[cookie.Value]))
	// log.Println("lobbyid", lobbyId)

	http.Redirect(w, r, "/uno/lobby/"+lobbyId, http.StatusSeeOther)
} // }}}

func serveLogout(w http.ResponseWriter, r *http.Request) { // {{{
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
} // }}}

func serveLobby(w http.ResponseWriter, r *http.Request) { // {{{
	cookie := checkForCookie(w, r)
	id, _ := strconv.Atoi(r.PathValue("id"))
	if lobby.LobbyCount > id {
		http.Redirect(w, r, "/uno", http.StatusSeeOther)
	}
	log.Printf("serving /uno/lobby/%v to %v", id, player.PlayerList[cookie.Value].Name)

	tmpl, err := template.ParseFiles("./web/static/uno/lobby.html")
	log.Println("Parsing...")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = tmpl.Execute(w, player.PlayerList[cookie.Value].Name)
	log.Println("Executing...")
	if err != nil {
		log.Fatal(err)
		return
	}
} // }}}

func randomKey(len int) (key []byte) { // {{{
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
} // }}}

func checkForCookie(w http.ResponseWriter, r *http.Request) *http.Cookie { // {{{
	cookie, err := r.Cookie("unoName")
	if err != nil || player.PlayerList[cookie.Value].Name == "" {
		http.Redirect(w, r, "/uno", http.StatusSeeOther)
	}
	return cookie
} // }}}
