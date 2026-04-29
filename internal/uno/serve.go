package uno

// Maybe use JWT
// Instead of storing auth as a

// TODO: Right now Everyone has complete access to refreshing and creating new accounts, we are only testing access tokens rn
// TODO: Use refresh auths to save as a cookie and only use small tokens
// TODO: Have user data saved even if refreshed via cookies
// TODO: Make a refresh auth that is tied to someone's name, save a database of names
// TODO: Implement in-band token refresh
// TODO: Hide the secret somewhere
// TODO: Allow ussers to create lobbies n shit
import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	snapws "github.com/Atheer-Ganayem/SnapWS"
	jwt "github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"

	serve "coollittlewebsite/pkg/serve"
)

type Envelope struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

// string of player is the username of the player
type State struct {
	Players map[string](Player)
}

type Player struct {
	Count int
}

var (
	globalCount int
	roomManager *snapws.RoomManager[string]
	games       map[*snapws.Room[string]](State)
	secret      []byte
)

// Main page and assets
func Serve(man **snapws.RoomManager[string]) {
	// Databases of players
	log.Println(sql.Drivers())
	// db, err := sql.Open("sqlite3", "./tmp/players.db")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// defer db.Close()
	// fmt.Println("Connected to the SQLite database successfully.")
	//
	// // Get the version of SQLite
	// var sqliteVersion string
	// err = db.QueryRow("select sqlite_version()").Scan(&sqliteVersion)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// Create Room Manager
	*man = snapws.NewRoomManager[string](nil)
	roomManager = *man
	roomManager.Upgrader.Use(checkAuth)

	roomManager.DefaultOnJoin = func(room *snapws.Room[string], conn *snapws.Conn, args ...any) {
		name, _ := snapws.GetArg[string](args, 0)
		log.Printf("%s Joined", name)

		// Give the user their name
		err := conn.SendJSON(
			context.Background(),
			Envelope{Type: "name", Data: map[string]string{"name": name}},
		)
		if err != nil {
			return
		}

		games[room].Players[name] = Player{Count: 0}
		updateAll(room)
	}

	roomManager.DefaultOnLeave = func(room *snapws.Room[string], conn *snapws.Conn, args ...any) {
		log.Println(room.Key)
		name, _ := snapws.GetArg[string](args, 0)
		log.Printf("%s Left", name)
		delete(games[room].Players, name)
		updateAll(room)
	}

	games = make(map[*snapws.Room[string]](State))

	// Create the "uno" room
	room := roomManager.Add("uno")
	games[room] = State{
		Players: make(map[string](Player)),
	}

	serve.ServeIndex("/uno")
	serve.ServeAssets("/uno")
	http.HandleFunc("GET /uno/auth", auth)
	http.HandleFunc("GET /uno/upgrade", upgrade)
	// http.HandleFunc("GET /uno/cookie", getCookie)
}

func upgrade(w http.ResponseWriter, r *http.Request) {
	log.Printf("upgrading %s", "/uno/upgrade")

	auth := r.URL.Query().Get("auth")
	roomName := r.URL.Query().Get("room")
	log.Printf("Validating query: '%s'", auth)

	token, err := jwt.Parse(
		auth,
		func(token *jwt.Token) (any, error) { return secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return
	}
	name, err := token.Claims.(jwt.MapClaims).GetSubject()
	if err != nil {
		return
	}

	room := roomManager.Get(roomName)
	conn, _, err := roomManager.Connect(w, r, roomName, name)
	if err != nil {
		return
	}

	defer conn.Close()

	for {
		env := Envelope{}
		err := conn.ReadJSON(&env)

		if snapws.IsFatalErr(err) {
			return // Connection closed
		} else if err != nil {
			fmt.Println("Non-fatal error:", err)
			continue
		}
		switch env.Type {
		case "increment":
			increment(name, room)
			updateAll(room)
		}
	}
}

// Generate a new key for the user
func auth(w http.ResponseWriter, r *http.Request) {
	// Generate a new key
	username := fmt.Sprintf("anonymous_%d", globalCount)
	globalCount += 1

	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Subject:   username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		},
	)
	s, err := t.SignedString(secret)
	if err != nil {
		return
	}

	w.Write([]byte(s))
}

// Generate a new key for the user
func checkAuth(w http.ResponseWriter, r *http.Request) error {
	auth := r.URL.Query().Get("auth")
	log.Println("Validating query", auth)

	if auth == "" {
		return snapws.NewMiddlewareErr(http.StatusBadRequest, "Empty auth query")
	}

	token, err := jwt.Parse(
		auth,
		func(token *jwt.Token) (any, error) { return secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	switch {
	case token.Valid:
		return nil
	case errors.Is(err, jwt.ErrTokenMalformed):
		return snapws.NewMiddlewareErr(http.StatusBadRequest, "Not a token")
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return snapws.NewMiddlewareErr(http.StatusBadRequest, "Invalid signature")
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return snapws.NewMiddlewareErr(
			http.StatusBadRequest,
			"Token is either expired or not active yet",
		)
	default:
		return snapws.NewMiddlewareErr(http.StatusBadRequest, "Couldn't Handle Token")
	}

	// cookie, err := r.Cookie(cookieName)
	// if err == http.ErrNoCookie {
	// 	log.Println("Cookie DENIED,  No Cookie")
	// 	http.Redirect(w, r, "/uno/cookie", http.StatusTemporaryRedirect)
	// 	return snapws.NewMiddlewareErr(http.StatusBadRequest, "You need a cookie")
	// } else if _, ok := names[cookie.Value]; !ok {
	// 	http.SetCookie(w, &http.Cookie{Name: cookieName, MaxAge: -1})
	// 	log.Println("Cookie DENIED,  Wrong Cookie")
	// 	return snapws.NewMiddlewareErr(http.StatusBadRequest, "Cookie is invalid")
	// }
	// return nil
}

// This increments the player associated with a conn in the gamestate and also sends a message to all clients to update their gamestate
func increment(name string, room *snapws.Room[string]) {
	player, ok := games[room].Players[name]
	if !ok {
		return
	}

	player.Count += 1
	games[room].Players[name] = player
}

func updateAll(room *snapws.Room[string]) {
	_, err := room.BroadcastJSON(
		context.Background(),
		Envelope{Type: "state", Data: games[room]},
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
}

// HTTP Middleware that blocks your connection if
// func checkCookie(w http.ResponseWriter, r *http.Request) error {
// 	log.Println("Checking cookie")
// 	cookie, err := r.Cookie(cookieName)
// 	if err == http.ErrNoCookie {
// 		log.Println("Cookie DENIED,  No Cookie")
// 		http.Redirect(w, r, "/uno/cookie", http.StatusTemporaryRedirect)
// 		return snapws.NewMiddlewareErr(http.StatusBadRequest, "You need a cookie")
// 	} else if _, ok := names[cookie.Value]; !ok {
// 		http.SetCookie(w, &http.Cookie{Name: cookieName, MaxAge: -1})
// 		log.Println("Cookie DENIED,  Wrong Cookie")
// 		return snapws.NewMiddlewareErr(http.StatusBadRequest, "Cookie is invalid")
// 	}
// 	return nil
// }
//
// func getCookie(w http.ResponseWriter, r *http.Request) {
// 	log.Println("Forcing Cookie")
// 	cookie, err := r.Cookie(cookieName)
// 	if err != http.ErrNoCookie {
// 		if _, ok := names[cookie.Value]; ok {
// 			return
// 		}
// 	}
//
// 	// Generate a new key
// 	username := fmt.Sprintf("anonymous_%d", globalCount)
// 	key := username + "_key"
// 	globalCount += 1
// 	names[key] = username
// 	state.Players[username] = Player{Count: 0}
//
// 	new_cookie := http.Cookie{
// 		Name:   cookieName,
// 		Value:  key,
// 		Path:   "/uno",
// 		MaxAge: 3600,
// 	}
// 	http.SetCookie(w, &new_cookie)
// }
