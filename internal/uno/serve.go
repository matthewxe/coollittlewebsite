package uno

import (
	"coollittlewebsite/pkg/serve"
)

// Main page and assets
func Serve() {
	serve.ServeIndex("/uno")
}
