package webhooks

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var GITHUB_WEBHOOK_SECRET string = os.Getenv("GITHUB_WEBHOOK_SECRET")

func GithubWebhookHTTP(w http.ResponseWriter, r *http.Request) {
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
	if !VerifySignature(signature, payload) {
		http.Error(w, "Signature verification failed", http.StatusForbidden)
		return
	}
	if err := GitPull(); err != nil {
		http.Error(w, "Failed to execute git pull", http.StatusInternalServerError)
		return
	}

	// Respond to GitHub with a success status
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received successfully"))
}

func VerifySignature(signature string, payload []byte) bool {
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

func GitPull() error {
	cmd := exec.Command("git", "pull")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("git pull failed: %s", string(output))
	}
	return nil
}
