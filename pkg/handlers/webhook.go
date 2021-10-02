package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func errorResponse(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func getSecret() string {
	return os.Getenv("SECRET_KEY")
}

func isValidSignature(secret string, signature string, body []byte) (bool, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write(body)
	if err != nil {
		return false, err
	}

	digest := mac.Sum(nil)
	return hex.EncodeToString(digest) == signature, nil
}

func handleWebhook(body []byte, signature string, secret string) (int, error) {
	validSignature, err := isValidSignature(secret, signature, body)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to check if signature is valid: %w+", err)
	}

	if !validSignature {
		return http.StatusBadRequest, fmt.Errorf("invalid signature")
	}

	return http.StatusOK, nil
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	secret := getSecret()
	if secret == "" {
		errorResponse(w, r, fmt.Errorf("SECRET_KEY has not been set"), http.StatusInternalServerError)
		return
	}

	signature := strings.TrimLeft(r.Header.Get("X-Hub-Signature-256"), "sha256=")
	if signature == "" {
		errorResponse(w, r, fmt.Errorf("signature header not found"), http.StatusBadRequest)
		return
	}

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, r, fmt.Errorf("failed to read body: %w+", err), http.StatusInternalServerError)
		return
	}

	statusCode, err := handleWebhook(bs, signature, secret)
	if err != nil {
		errorResponse(w, r, err, statusCode)
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": "Success"})
}
