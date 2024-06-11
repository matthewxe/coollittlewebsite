package uno

import (
	"bytes"
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
	send map[int]chan []byte
} //}

func newPlayer(name string) (*Player, string) {
	newplayer := &Player{lobby: make(map[int]*Lobby), send: make(map[int]chan []byte), Name: name, conn: make(map[int]*websocket.Conn)}
	key := randomKey(24)
	return newplayer, key
}

func randomKey(len int) string { //{
	var key []byte
	for i := 0; i < len; i++ {
		excluded := []int{1, 26, 59}
		random := randIntExclude(93, excluded)
		key = append(key, byte(random+33))
	}
	return string(key)
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
	space   = []byte{' '}
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
	player.conn[id].SetPongHandler(func(string) error { player.conn[id].SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := player.conn[id].ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Printf("Recieved message from %s -> lobby %v with the message %s", player.Name, id, message)
		player.lobby[id].broadcast <- message
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
			log.Printf("Recieved message from lobby %v -> %s with the message %s", id, player.Name, message)
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
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(player.send[id])
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-player.send[id])
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			log.Printf("%s got ticked off in lobby %v", player.Name, id)
			player.conn[id].SetWriteDeadline(time.Now().Add(writeWait))
			if err := player.conn[id].WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
} //}

// vim:foldmethod=marker:foldmarker=//{,//}:
