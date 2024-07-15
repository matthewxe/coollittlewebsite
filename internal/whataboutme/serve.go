package whataboutme

import (
	"net/http"

	"coollittlewebsite/pkg/serve"
)

func Serve() {
	serve.ServeIndex("/whataboutme")
	http.HandleFunc("GET /whataboutme/", serve.ServeAssets)
}
