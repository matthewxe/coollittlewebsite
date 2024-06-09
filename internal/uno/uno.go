package uno

import (
	"coollittlewebsite/internal/serve/assets"
	"fmt"
	"html/template"
	"log"
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

	// Logging in and out
	http.HandleFunc("GET /uno/login", serveLogin)
	http.HandleFunc("GET /uno/logout", serveLogout)

	// TODO: Creating a lobby
	http.HandleFunc("GET /uno/create", serveCreate)
	http.HandleFunc("GET /uno/list", serveList)

	// TODO: Serve a lobby
	http.HandleFunc("GET /uno/lobby/{id}", serveLobby)
	http.HandleFunc("GET /uno/lobby/{id}/ws", serveLobbyWs)
}

func serveIndex(w http.ResponseWriter, r *http.Request) { //{
	cookie := checkForCookie(w, r)
	if cookie == nil {
		return
	}
	log.Println("serving /uno")

	tmpl, err := template.ParseFiles("./web/static/uno/index.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	tmpl_err := tmpl.Execute(w, playerList[cookie.Value].Name)
	if tmpl_err != nil {
		log.Fatal(err)
		return
	}
} //}

func serveList(w http.ResponseWriter, r *http.Request) { //{
	checkForCookie(w, r)
	var ready string
	var ongoing string
	var done string
	for i, lobbi := range lobbyList {
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
} //}

func serveLogin(w http.ResponseWriter, r *http.Request) { //{
	log.Println("serving /uno/login")
	cookie, err := r.Cookie("unoName")
	if err != nil || playerList[cookie.Value] == nil || playerList[cookie.Value].Name == "" {
		err := r.ParseForm()
		if err != nil {
			log.Fatal("Error to parse form")
		} else if r.Form.Get("name") != "" {
			newplayer, key := newPlayer(r.Form.Get("name"))
			playerList[key] = newplayer
			cookieNew := &http.Cookie{}
			cookieNew.Name = "unoName"
			cookieNew.Value = key
			cookieNew.Expires = time.Now().Add(365 * 24 * time.Hour) // After 1 year
			cookieNew.Secure = true
			// cookieNew.Secure = false
			cookieNew.HttpOnly = true
			cookieNew.Path = "/uno"
			http.SetCookie(w, cookieNew)
			cookie = cookieNew
		} else {
			http.ServeFile(w, r, "./web/static/uno/name.html")
			return
		}
	}
	http.Redirect(w, r, "/uno", http.StatusSeeOther)
} //}

func serveLogout(w http.ResponseWriter, r *http.Request) { //{
	cookie := checkForCookie(w, r)
	if cookie == nil {
		return
	}
	log.Println("serving /uno/logout to \"", cookie.Value, "\"")

	delete(playerList, cookie.Value)
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
} //}

func serveCreate(w http.ResponseWriter, r *http.Request) { //{
	cookie := checkForCookie(w, r)
	if cookie == nil {
		return
	}
	log.Println("serving /uno/create to ", playerList[cookie.Value].Name)

	lobby := newLobby(playerList[cookie.Value])
	go lobby.run()
	log.Printf("creating lobby %d", lobby.Id)

	http.Redirect(w, r, fmt.Sprintf("/uno/lobby/%v", lobby.Id), http.StatusSeeOther)
} //}

func serveLobby(w http.ResponseWriter, r *http.Request) { //{
	cookie := checkForCookie(w, r)
	if cookie == nil {
		return
	}
	id, _ := strconv.Atoi(r.PathValue("id"))
	if lobbyCount <= id {
		log.Println("invalid lobby")
		http.Redirect(w, r, "/uno", http.StatusSeeOther)
		return
	}
	lobber := lobbyList[id]
	if lobber.State != 0 {
		log.Println("invalid lobby")
		http.Redirect(w, r, "/uno", http.StatusSeeOther)
		return
	}

	log.Printf("serving /uno/lobby/%v to %v", id, playerList[cookie.Value].Name)

	tmpl, err := template.ParseFiles("./web/static/uno/lobby.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = tmpl.Execute(w, lobber)
	if err != nil {
		log.Fatal(err)
		return
	}
} //}

func serveLobbyWs(w http.ResponseWriter, r *http.Request) { //{
	cookie := checkForCookie(w, r)
	if cookie == nil {
		return
	}
	id, _ := strconv.Atoi(r.PathValue("id"))
	if lobbyCount <= id {
		log.Println("invalid lobby")
		http.Redirect(w, r, "/uno", http.StatusSeeOther)
		return
	}
	lobber := lobbyList[id]
	if lobber.State != 0 {
		log.Println("invalid lobby")
		http.Redirect(w, r, "/uno", http.StatusSeeOther)
		return
	}
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	player := playerList[cookie.Value]
	player.lobby[id] = lobber
	player.conn[id] = conn
	player.send[id] = make(chan []byte, 256)
	player.lobby[id].register <- player

	go player.writePump(id)
	go player.readPump(id)
} //}

func checkForCookie(w http.ResponseWriter, r *http.Request) *http.Cookie { //{
	cookie, err := r.Cookie("unoName")
	if err != nil || playerList[cookie.Value] == nil || playerList[cookie.Value].Name == "" {
		http.Redirect(w, r, "/uno/login", http.StatusSeeOther)
		return nil
	}
	return cookie
} //}
// vim:foldmethod=marker:foldmarker=//{,//}:
