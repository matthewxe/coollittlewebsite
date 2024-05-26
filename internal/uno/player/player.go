package player

// import "coollittlewebsite/internal/uno/lobby"

var PlayerList = make(map[string]Player)

type Player struct {
	// hub []*lobby.Lobby

	// The websocket connection.
	// conn *lobby..websocket.Conn

	Name string
}
