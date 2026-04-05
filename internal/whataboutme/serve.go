package whataboutme

import "coollittlewebsite/pkg/serve"

func Serve() {
	serve.ServeIndex("/whataboutme")
	serve.ServeAssets("/whataboutme")
}
