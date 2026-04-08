package uno

// "encoding/json"
import (
	"context"
	"fmt"
	"log"
	"net/http"

	snapws "github.com/Atheer-Ganayem/SnapWS"

	serve "coollittlewebsite/pkg/serve"
)

var (
	cookieName  = "uno"
	globalCount int
	// key -> username
	keys     map[string](string)
	manager  *snapws.Manager[string]
	state    State
	upgrader *snapws.Upgrader
)

// Main page and assets
func Serve(man **snapws.Manager[string]) {
	upgrader = snapws.NewUpgrader(nil)
	upgrader.Use(checkCookie)
	// TODO: Auth middleware?
	*man = snapws.NewManager[string](upgrader)
	manager = *man
	// manager.Shutdown()

	state = State{}
	state.Players = make(map[string](Player))
	keys = make(map[string](string))

	// Hooks that do an action whenever someone connects
	// manager.OnRegister = onRegister
	// manager.OnUnregister = onUnregister

	// Serve the /uno index and js
	serve.ServeIndex("/uno")
	serve.ServeAssets("/uno")
	// Upgrade the websocket
	http.HandleFunc("GET /uno/upgrade", handler)
	http.HandleFunc("GET /uno/cookie", forceCookies)
}

type Envelope struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// string of player is the username of the player
type State struct {
	Players map[string](Player)
}

type Player struct {
	Count int
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("upgrading %s", "/uno/upgrade")

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return
	}
	key := cookie.Value

	if c := manager.Get(key); c != nil {
		c.CloseWithCode(1, "Opened in a different window")
	}

	conn, err := manager.Connect(key, w, r)
	if err != nil {
		return
	}

	username := keys[key]

	defer conn.Close()
	defer delete(keys, username)
	defer delete(state.Players, conn.Key)
	defer fmt.Println("Closed", key)

	// Give the user their name
	err = conn.SendJSON(
		context.Background(),
		Envelope{Type: "name", Data: map[string]string{"name": username}},
	)
	if err != nil {
		return
	}

	updateAll()

	for {
		// var msg string
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
			increment(conn)
		}
	}
}

// This increments the conn and also sends a message to the client to update
// NOTE: WE DO NOT WANT THIS BECAUSE IT DOES NOT UPDATE THE ENTIRE GAMESTATE WHICH WE WANT INSTEAD
// only for testing
func increment(conn *snapws.ManagedConn[string]) {
	player, ok := state.Players[keys[conn.Key]]
	if !ok {
		return
	}

	player.Count += 1
	state.Players[keys[conn.Key]] = player
	// err := conn.SendJSON(
	// 	context.Background(),
	// 	Envelope{
	// 		Type: "count", Data: map[string]int{"count": player.Count},
	// 	},
	// )
	// if err != nil {
	// 	return
	// }
	updateAll()
}

// This sends the client an updated gamestate which they will rebuild their structure from
func update(conn *snapws.ManagedConn[string]) {
	err := conn.SendJSON(
		context.Background(),
		Envelope{
			Type: "state", Data: state,
		},
	)
	if err != nil {
		return
	}
}

func updateAll() {
	// Increment everyone when someone joins
	for _, c := range manager.GetAllConns() {
		update(c)
	}
}

func checkCookie(w http.ResponseWriter, r *http.Request) error {
	log.Println("yo")
	cookie, err := r.Cookie(cookieName)
	if err == http.ErrNoCookie {
		return snapws.NewMiddlewareErr(http.StatusBadRequest, "You need a cookie")
	} else if _, ok := keys[cookie.Value]; !ok {
		http.SetCookie(w, &http.Cookie{Name: cookieName, MaxAge: -1})
		return snapws.NewMiddlewareErr(http.StatusBadRequest, "Cookie is invalid")
	}
	return nil
}

func forceCookies(w http.ResponseWriter, r *http.Request) {
	log.Println("yo")
	cookie, err := r.Cookie(cookieName)
	if err != http.ErrNoCookie {
		if _, ok := keys[cookie.Value]; ok {
			return
		}
	}

	// Generate a new key
	username := fmt.Sprintf("anonymous_%d", globalCount)
	key := username + "_key"
	globalCount += 1
	keys[key] = username
	state.Players[username] = Player{Count: 0}

	new_cookie := http.Cookie{
		Name:   cookieName,
		Value:  key,
		Path:   "/uno",
		MaxAge: 3600,
	}
	http.SetCookie(w, &new_cookie)
}

// This is some dummy hooks.
// In real world you might send a message to update the user's status for the other connected users.
// func onRegister(conn *snapws.ManagedConn[string]) {
// 	id := conn.Key
// 	manager.BroadcastString(context.TODO(), []byte(id+" is online!"), id)
// }
//
// func onUnregister(conn *snapws.ManagedConn[string]) {
// 	id := conn.Key
// 	conn.Manager.BroadcastString(context.TODO(), []byte(id+" is offline"), id)
// }
