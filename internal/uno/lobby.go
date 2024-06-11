package uno

import (
	"log"
)

var lobbyList []*Lobby
var lobbyCount int = 0

type Lobby struct { //{
	Id int

	// Leader also exists in the Players map
	Leader *Player

	// This if a *Player exists in a map it means it is inside a lobby so it
	// Has joined the lobby
	// But true means that ists connected by websocket, and false means no
	Players map[*Player]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	logs [][]byte
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
	lobby.Players[leader] = false
	lobbyCount++
	lobbyList = append(lobbyList, lobby)
	log.Println(lobbyList)

	leader.send[lobby.Id] = make(chan []byte, 256)
	return lobby
} //}

func (lobby *Lobby) run() { //{
	for {
		select {
		case player := <-lobby.register:
			// log.Printf("%s registered [lobby %v]", player.Name, lobby.Id)
			for _, v := range lobby.logs {
				player.send[lobby.Id] <- v
				log.Printf("%s", v)
			}
			lobby.Players[player] = true
		case player := <-lobby.unregister:
			if _, ok := lobby.Players[player]; ok {
				// delete(h.Players, client)
				// log.Printf("%s unregistered [lobby %v]", player.Name, lobby.Id)
				lobby.Players[player] = false
				if _, ok := player.send[lobby.Id]; !ok {
					close(player.send[lobby.Id])
				}
			}
		case message := <-lobby.broadcast:
			lobby.logs = append(lobby.logs, message)
			for player := range lobby.Players {
				// log.Printf("Lobby.broadcast list players: %s", player.Name)
				select {
				case player.send[lobby.Id] <- message:
					// log.Printf("%s messaged [lobby %v]: '%s'", player.Name, lobby.Id, message)
				default:
					// delete(lobby.players, player)
					// log.Printf("%s failed to message and unregistered [lobby %v]", player.Name, lobby.Id)
					lobby.Players[player] = false
					if _, ok := player.send[lobby.Id]; !ok {
						close(player.send[lobby.Id])
					}
				}
			}
		}
	}
} //}
// vim:foldmethod=marker:foldmarker=//{,//}:
