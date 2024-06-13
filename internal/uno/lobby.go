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
	logs []ChatJSON

	// Inbound messages from the clients.
	broadcast chan JSON

	// Register requests from the clients.
	register chan *Player

	// Unregister requests from clients.
	unregister chan *Player

	// 0 Not Playing
	// 1 Playing
	// 2 Game finished
	State int
} //}

func newLobby(leader *Player) *Lobby { //{
	var lobby = &Lobby{
		Players:    make(map[*Player]bool),
		broadcast:  make(chan JSON),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		Leader:     leader,
		State:      0,
		Id:         lobbyCount,
	}
	lobby.Players[leader] = false
	lobbyCount++
	lobbyList = append(lobbyList, lobby)
	// log.Println(lobbyList)

	leader.send[lobby.Id] = make(chan JSON, 256)
	return lobby
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

type ChatJSON struct { //{
	Type   string
	Player string
	Text   string
	Date   int64
} //}

type ChatLogJSON struct { //{
	Type string
	Log  []ChatJSON
} //}

type StartJSON struct { //{
	Type string
} //}

type JSON struct {
	Type   string
	Player *Player
	JSON   []byte
}

func (l Lobby) MessageLog() JSON { // //{
	messagelog := ChatLogJSON{"messagelog", nil}

	messagelog.Log = append(messagelog.Log, l.logs...)

	log.Println("SavedMessageLog", l.logs)
	log.Println("MessageLog", messagelog.Log)
	marshal, err := json.Marshal(messagelog)
	if err != nil {
		return JSON{"error", nil, nil}
	}
	return JSON{"messagelog", nil, marshal}
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
		player.send[l.Id] <- JSON{"message", nil, l.JsonifyBytify(player)}
	}
} // //}

func (lobby *Lobby) run() { //{
	for {
		select {
		case player := <-lobby.register:
			messagelog := lobby.MessageLog()
			if messagelog.Type != "error" {
				player.send[lobby.Id] <- messagelog
			}
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
			switch message.Type {
			case "message":
				var chat ChatJSON
				if err := json.Unmarshal(message.JSON, &chat); err == nil {
					lobby.logs = append(lobby.logs, chat)
				}
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
			case "start":
			default:
				return
			}
		}
	}
} //}
// vim:foldmethod=marker:foldmarker=//{,//}:
