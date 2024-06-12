package uno

import (
	"encoding/json"
	"log"
)

var lobbyList []*Lobby
var lobbyCount int = 0

type Lobby struct { //{
	// Id in the lobbyList
	Id int

	// Leader also exists in the Players map
	Leader *Player

	// This if a *Player exists in a map it means it is inside a lobby
	// True means its currently connected to the websocket
	// False means its currently not connected to the websocket
	Players map[*Player]bool

	// Logs of chat messages
	logs []MessageJSON

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Player

	// Unregister requests from clients.
	unregister chan *Player

	// 0 Not Playing
	// 1 Playing
	// 2 Game finished
	State int
} //}

type PlayerJSON struct { //{
	Type    string
	Name    string
	Active  bool
	Current bool
} //}

type LobbyJSON struct { //{
	Type    string
	Id      int
	Leader  PlayerJSON
	Players []PlayerJSON
	State   int
} //}

type MessageJSON struct { //{
	Type   string
	Player string
	Text   string
	Date   int
} //}

type MessageLogJSON struct { //{
	Type string
	Log  []MessageJSON
} //}

func (l Lobby) MessageLog() []byte { // //{
	log := MessageLogJSON{"messagelog", nil}

	log.Log = append(log.Log, l.logs...)

	marshal, err := json.Marshal(log)
	if err != nil {
		return nil
	}
	return marshal
} // //}

func (l Lobby) Jsonify(p *Player) LobbyJSON { //{
	leader := l.Leader
	var playerlist []PlayerJSON
	for player := range l.Players {
		if player != leader {
			playerlist = append(playerlist, PlayerJSON{"player", player.Name, l.Players[player], player == p})
		}
	}
	return LobbyJSON{"status", l.Id, PlayerJSON{"player", leader.Name, l.Players[leader], leader == p}, playerlist, l.State}
} //}

func (l Lobby) JsonifyBytify(p *Player) []byte { //{
	marshal, err := json.Marshal(l.Jsonify(p))
	if err != nil {
		return nil
	}
	return marshal
} //}

// Sends a JSON to the all players in a lobby to tell them to update their
// lobby status so its always up to date
func (l Lobby) UpdatePlayers() { // //{
	for player := range l.Players {
		player.send[l.Id] <- l.JsonifyBytify(player)
	}
} // //}

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
			player.send[lobby.Id] <- lobby.MessageLog()
			lobby.Players[player] = true
			lobby.UpdatePlayers()
		case player := <-lobby.unregister:
			if _, ok := lobby.Players[player]; ok {
				lobby.Players[player] = false
				lobby.UpdatePlayers()
				if _, ok := player.send[lobby.Id]; !ok {
					close(player.send[lobby.Id])
				}
			}
		case message := <-lobby.broadcast:
			// lobby.logs = append(lobby.logs, message)
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
