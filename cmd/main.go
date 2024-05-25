package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os/exec"

	_ "modernc.org/sqlite"
	// "text/template"
)

const ABOUT_DIR string = "/whataboutme"
const GITHUB_WEBHOOK_DIR string = "/heyyyyhaveyouheardaboutthisthingcalledahook" // Temporary solution for live reload over github
const STATIC_DIR string = "web/static"
const PORT string = ":8080"

var GITHUB_WEBHOOK_SECRET string = os.Getenv("GITHUB_WEBHOOK_SECRET")
var api_key = random_key(32)
var db *sql.DB

func main() {
	db, err := sql.Open("sqlite", "./data/data.db")
	if err != nil {
		return
	}
	defer db.Close()
	setup_sql(db)

	log.Printf("The key is\n%s", api_key)
	os.WriteFile("/tmp/coollittlewebsite_key", []byte(api_key), 0666)
	log.Print("Listening on " + PORT + "...")

	http.HandleFunc("GET "+ABOUT_DIR,
		func(w http.ResponseWriter, r *http.Request) {
			log.Print("serving about me")
			http.ServeFile(w, r, STATIC_DIR+"/index.html")
		})
	http.HandleFunc("GET "+ABOUT_DIR+"/", serve_assets)
	http.HandleFunc("GET "+ABOUT_DIR+"/blog/{id}", serve_blog)
	http.HandleFunc("GET "+ABOUT_DIR+"/addanewpostyoubingus", serve_new_post)
	http.HandleFunc("POST "+ABOUT_DIR+"/addanewpostyoubingus", serve_new_post)
	http.HandleFunc("POST "+GITHUB_WEBHOOK_DIR, github_webhook)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func random_key(len int) (key []byte) {
	for i := 0; i < len; i++ {
		key = append(key, byte((rand.Int()%97)+33))
	}
	return key
}

func setup_sql(db *sql.DB) {
	table_create := `CREATE TABLE blogs (
		id INTEGER PRIMARY KEY,
		thumbnail BLOB,
		content TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)`
	db.Exec(table_create)
}

func serve_assets(w http.ResponseWriter, r *http.Request) {
	var page string = strings.TrimPrefix(r.RequestURI, ABOUT_DIR)
	if page == "/" {
		http.Redirect(w, r, ABOUT_DIR, http.StatusPermanentRedirect)
	}
	page = STATIC_DIR + page
	log.Print("serving asset " + page)

	_, err := os.Stat(page)
	if os.IsNotExist(err) {
		log.Print("failed to serve asset " + page)
		http.Redirect(w, r, ABOUT_DIR, http.StatusPermanentRedirect)
		return
	}
	http.ServeFile(w, r, page)
}

func serve_new_post(w http.ResponseWriter, r *http.Request) {
	log.Print("add a new post")
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, "web/template/newpostprompt.html")
	case http.MethodPost:
		r.ParseForm()

		log.Printf("'%s' != '%s' ? %b", string(api_key), r.Form.Get("key"), string(api_key) != string(r.Form.Get("key")))

		if string(api_key) != string(r.Form.Get("key")) {
			log.Print("fail")
			http.Redirect(w, r, ABOUT_DIR, http.StatusFound)
			return
		}
		log.Print("fail")
		w.Write([]byte("Oogabooga"))
	}
}
func serve_blog(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
	id := r.FormValue("id")

	rows, err := db.Query("SELECT content from blogs WHERE id = id", id)
	if err != nil {
		return
	}
	var content string
	rows.Scan(content)

	w.Write([]byte(content))
}

func github_webhook(w http.ResponseWriter, r *http.Request) {
	log.Print("Github webhook activated")

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	signature := r.Header.Get("X-Hub-Signature")
	if signature == "" {
		http.Error(w, "Signature missing", http.StatusBadRequest)
		return
	}
	if !verifySignature(signature, payload) {
		http.Error(w, "Signature verification failed", http.StatusForbidden)
		return
	}
	if err := gitPull(); err != nil {
		http.Error(w, "Failed to execute git pull", http.StatusInternalServerError)
		return
	}

	// Respond to GitHub with a success status
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received successfully"))
	//         if self.path == '/github-webhook/':
	//             print("Handling GitHub webhook request")  # Debugging line
	//             # Read the payload from the request
	//             content_length = int(self.headers['Content-Length'])
	//             payload = self.rfile.read(content_length)
	//
	//             # Verify the webhook signature (if using a secret)
	//             if verify_github_webhook_signature(self.headers, payload, "DATAPACK"):
	//                 run_scripts(payload, self, DATAPACK)
	//             elif verify_github_webhook_signature(self.headers, payload, "RESOURCEPACK"):
	//                 run_scripts(payload, self, RESOURCEPACK)
	//             else:
	//                 self.send_response(HTTPStatus.UNAUTHORIZED)
	//                 self.end_headers()
	//                 self.wfile.write(b'Unauthorized')

}

func verifySignature(signature string, payload []byte) bool {
	// GitHub sends the signature in the format "sha1=XXXXXXXXX"
	parts := strings.SplitN(signature, "=", 2)
	if len(parts) != 2 || parts[0] != "sha1" {
		return false
	}

	// Compute the HMAC digest of the payload using the secret
	mac := hmac.New(sha1.New, []byte(GITHUB_WEBHOOK_SECRET))
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	// Compare the computed digest with the signature
	return hmac.Equal([]byte(parts[1]), []byte(expectedMAC))
}

func gitPull() error {
	cmd := exec.Command("git", "pull")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("git pull failed: %s", string(output))
	}
	return nil
}
