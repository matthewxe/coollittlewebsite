package uno

import "log"

var lobbyList []*Lobby
var lobbyCount int = 0

type Lobby struct { //{
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
} //}

func newLobby(leader *Player) *Lobby { //{
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
	lobbyCount++
	lobbyList = append(lobbyList, lobby)
	log.Println(lobbyList)
	return lobby
} //}

func (h *Lobby) run() { //{
	for {
		select {
		case client := <-h.register:
			h.Players[client] = true
		case client := <-h.unregister:
			if _, ok := h.Players[client]; ok {
				delete(h.Players, client)
				close(client.send[h.Id])
			}
		case message := <-h.broadcast:
			for client := range h.Players {
				select {
				case client.send[h.Id] <- message:
				default:
					close(client.send[h.Id])
					delete(h.Players, client)
				}
			}
		}
	}
} //}
// vim:foldmethod=marker:foldmarker=//{,//}:
