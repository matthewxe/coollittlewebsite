package uno

import "log"

var lobbyList []*Lobby
var lobbyCount int = 0

type Lobby struct {
	Id int

	Leader *Player

	Players map[*Player]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Player

	// Unregister requests from clients.
	unregister chan *Player

	// If the
	State int
	// 0 Not Playing
	// 1 Playing
	// 2 Game finished
}

func newLobby(leader *Player) (id int) {
	var lobby = &Lobby{
		Players:    make(map[*Player]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		Leader:     leader,
		State:      0,
		Id:         lobbyCount,
	}
	lobby.Players[leader] = true
	id = lobbyCount
	lobbyCount++
	lobbyList = append(lobbyList, lobby)
	log.Println(lobbyList)
	return id
}

func (h *Lobby) run() {
	for {
		select {
		case client := <-h.register:
			h.Players[client] = true
		case client := <-h.unregister:
			if _, ok := h.Players[client]; ok {
				delete(h.Players, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.Players {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.Players, client)
				}
			}
		}
	}
}
//vi:foldmethod=marker:foldmarker={//,}//:
