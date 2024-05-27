package lobby

// "bytes"
// "log"
// "net/http"
// "time"
// "github.com/gorilla/websocket"
import (
	"coollittlewebsite/internal/uno/player"
	"log"
)

var LobbyList []Lobby
var LobbyCount int = 0

type Lobby struct {
	Id int

	Leader player.Player

	Players map[player.Player]bool

	// If the
	State int
	// 0 Not Playing
	// 1 Playing
	// 2 Game finished
}

func NewLobby(leader player.Player) (id int) {
	var lobby = Lobby{
		Players: make(map[player.Player]bool),
		Leader:  leader,
		State:   0,
		Id:      LobbyCount,
	}
	id = LobbyCount
	LobbyCount++
	LobbyList = append(LobbyList, lobby)
	log.Println(LobbyList)
	return id
}
