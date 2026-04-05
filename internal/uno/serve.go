package uno

import (
	"context"
	"fmt"
	"log"
	"net/http"

	snapws "github.com/Atheer-Ganayem/SnapWS"

	serve "coollittlewebsite/pkg/serve"
)

var upgrader *snapws.Upgrader

// Main page and assets
func Serve() {
	upgrader = snapws.NewUpgrader(nil)

	serve.ServeIndex("/uno")
	serve.ServeAssets("/uno")

	// Upgrade the websocket
	http.HandleFunc("GET /uno/upgrade", handler)
	serve.RedirectSlash("/uno/upgrade")
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("upgrading %s", "/uno/upgrade")

	conn, err := upgrader.Upgrade(w, r)
	if err != nil {
		return
	}
	defer fmt.Println("Closed")
	defer conn.Close()

	for {
		fmt.Println("Successful upgrade!")
		data, err := conn.ReadString()
		if snapws.IsFatalErr(err) {
			return // Connection closed
		} else if err != nil {
			fmt.Println("Non-fatal error:", err)
			continue
		}

		err = conn.SendString(context.TODO(), data)
		if snapws.IsFatalErr(err) {
			return // Connection closed
		} else if err != nil {
			fmt.Println("Non-fatal error:", err)
			continue
		}
	}
}
