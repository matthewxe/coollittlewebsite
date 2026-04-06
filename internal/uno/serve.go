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
	upgrader    *snapws.Upgrader
	manager     *snapws.Manager[string]
	globalCount int
	state       State
)

// Main page and assets
func Serve() {
	upgrader = snapws.NewUpgrader(nil)
	// TODO: Auth middleware
	manager = snapws.NewManager[string](upgrader)
	// defer manager.Shutdown()

	state = State{}
	state.Players = make(map[string](Player))

	// Hooks that do an action whenever someone connects
	// manager.OnRegister = onRegister
	// manager.OnUnregister = onUnregister

	// Serve the /uno index and js
	serve.ServeIndex("/uno")
	serve.ServeAssets("/uno")
	// Upgrade the websocket
	http.HandleFunc("GET /uno/upgrade", handler)
	serve.RedirectSlash("/uno/upgrade")
}

type Envelope struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// string of player is the key of a connection
type State struct {
	Players map[string](Player)
}

type Player struct {
	Name  string
	Count int
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("upgrading %s", "/uno/upgrade")

	name := fmt.Sprintf("anonymous_%d", globalCount)
	globalCount += 1

	conn, err := manager.Connect(name, w, r)
	if err != nil {
		return
	}

	state.Players[conn.Key] = Player{Name: name, Count: 0}

	defer fmt.Println("Closed", name)
	defer conn.Close()
	defer delete(state.Players, conn.Key)

	// Give the user their name
	err = conn.SendJSON(
		context.Background(),
		Envelope{Type: "name", Data: map[string]string{"name": name}},
	)
	if err != nil {
		return
	}

	// Increment everyone when someone joins
	for _, c := range manager.GetAllConns(conn.Key) {
		increment(c)
	}

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
	player, ok := state.Players[conn.Key]
	if !ok {
		return
	}

	player.Count += 1
	state.Players[conn.Key] = player
	err := conn.SendJSON(
		context.Background(),
		Envelope{
			Type: "count", Data: map[string]int{"count": player.Count},
		},
	)
	if err != nil {
		return
	}
}

// This sends the client an updated gamestate which they will rebuild their structure from
// func update(conn *snapws.ManagedConn[string]) {
// 	count[conn.Key] += 1
// 	err := conn.SendJSON(
// 		context.Background(),
// 		Envelope{
// 			Type: "state", Data: map[string]int{"count": count[conn.Key]},
// 		},
// 	)
// 	if err != nil {
// 		return
// 	}
// }

// this was inside the forloopl
// if msg[0] == 49 {
//
// }

//
// if targetConn := manager.Get(msg.To); targetConn != nil {
// 	rm := receivedMsg{Text: fmt.Sprintf("%s: %s", name, msg.Text), From: name}
// 	if err := targetConn.SendJSON(context.TODO(), rm); err != nil {
// 		fmt.Printf("error sending message from %s to %s: %v\n", name, msg.To, err)
// 	}
// }

// log.Println(manager.GetAllConns(""))
// data, err := conn.ReadString()
// if snapws.IsFatalErr(err) {
// 	return // Connection closed
// } else if err != nil {
// 	fmt.Println("Non-fatal error:", err)
// 	continue
// }
//
// err = conn.SendString(context.TODO(), data)
// if snapws.IsFatalErr(err) {
// 	return // Connection closed
// } else if err != nil {
// 	fmt.Println("Non-fatal error:", err)
// 	continue
// }

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
