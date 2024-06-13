package uno

import (
	// "encoding/json"
	"encoding/json"
	"html"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

var playerList = make(map[string]*Player)

type Player struct { //{
	Name string

	// Multiple lobbies
	lobby map[int]*Lobby

	// The websocket connection.
	conn map[int]*websocket.Conn

	// Buffered channel of outbound messages.
	// Formatted in JSON
	// Use DecodeJSON() to get the corresponding struct
	send map[int]chan JSON
} //}

func newPlayer(name string) (*Player, string) { //{
	newplayer := &Player{lobby: make(map[int]*Lobby),
		send: make(map[int]chan JSON),
		Name: html.EscapeString(name),
		conn: make(map[int]*websocket.Conn)}
	var len int = 24
	var key []byte
	for i := 0; i < len; i++ {
		// Only allow characters in a cookie-value
		excluded := []int{1, 26, 59}
		random := randIntExclude(93, excluded)
		// +33 aligns it to ASCII
		key = append(key, byte(random+33))
	}
	return newplayer, string(key)
} //}

func randIntExclude(top int, excluded []int) (random int) { //{
	random = (rand.Int() % top)
	for _, v := range excluded {
		if random == v {
			return randIntExclude(top, excluded)
		}
	}
	return
} //}

func MessageUnmarshal(message []byte) interface{} {
	var chat ChatJSON

	err := json.Unmarshal(message, &chat)
	if err == nil && chat.Type == "message" {
		return chat
	}

	var start StartJSON
	err = json.Unmarshal(message, &start)
	if err == nil && chat.Type == "start" {
		return start
	}
	return nil
}

const ( //{
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
} //}

// Player is a middleman between the websocket connection and the lobby.

// readPump pumps messages from the websocket connection to the lobby.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (player *Player) readPump(id int) { //{
	defer func() {
		player.lobby[id].unregister <- player
		player.conn[id].Close()
	}()
	player.conn[id].SetReadLimit(maxMessageSize)
	player.conn[id].SetReadDeadline(time.Now().Add(pongWait))
	player.conn[id].SetPongHandler(func(string) error {
		player.conn[id].SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := player.conn[id].ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// log.Printf("Recieved message from %s -> lobby %v with the message %s", player.Name, id, message)

		var Type string

		unmarsh := MessageUnmarshal(message)
		switch unmarshtype := unmarsh.(type) {
		case ChatJSON:
			Type = "message"
			unmarshtype.Date = time.Now().UnixMilli()
			unmarshtype.Player = player.Name
			unmarsh = unmarshtype
		case StartJSON:
			Type = "start"
		default:
			return
		}
		marsh, err := json.Marshal(unmarsh)
		if err != nil {
			log.Fatal("Schiesse")
			return
		}

		player.lobby[id].broadcast <- JSON{Type, player, marsh}
	}
} //}

// writePump pumps messages from the lobby to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (player *Player) writePump(id int) { //{
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		player.conn[id].Close()
	}()
	for {
		select {
		case message, ok := <-player.send[id]:
			// log.Printf("Recieved message from lobby %v -> %s with the message %s", id, player.Name, message)
			player.conn[id].SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The lobby closed the channel.
				player.conn[id].WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := player.conn[id].NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message.JSON)

			// Add queued chat messages to the current websocket message.
			n := len(player.send[id])
			for i := 1; i < n; i++ {
				w.Write(message.JSON)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// log.Printf("%s got ticked off in lobby %v", player.Name, id)
			player.conn[id].SetWriteDeadline(time.Now().Add(writeWait))
			if err := player.conn[id].WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
} //}

// vim:foldmethod=marker:foldmarker=//{,//}:
