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

	"github.com/itchyny/gojq"
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

type route struct {
	events []string
	name   string
	query  string
	route  string
}

func getRoutes(envVars []string) []route {
	routes := []route{}

	for _, e := range envVars {
		envParts := strings.SplitN(e, "=", 2)
		key := envParts[0]
		val := envParts[1]

		if strings.HasPrefix(key, "RULE_") {
			ruleParts := strings.SplitN(val, "|", 4)
			if len(ruleParts) != 3 {
				fmt.Printf("Failed to parse rule: %s, wrong number of parts\n", key)
				continue
			}

			routes = append(routes, route{
				events: strings.Split(ruleParts[0], ","),
				query:  ruleParts[2],
				name:   key,
				route:  ruleParts[1],
			})
		}
	}

	return routes
}

func routeWebhook(body []byte, route route) error {
	fmt.Printf("Sending webhook to %s\n", route.route)
	return nil
}

func stringArrayContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func routeMatches(body []byte, event string, route route) (bool, error) {
	query, err := gojq.Parse(route.query)
	if err != nil {
		return false, err
	}

	if !stringArrayContains(route.events, event) {
		return false, nil
	}

	input := map[string]interface{}{}
	err = json.Unmarshal(body, &input)
	if err != nil {
		return false, err
	}

	iter := query.Run(input)
	for {
		v, ok := iter.Next()
		if !ok {
			// Finished iteration, must have found a match
			break
		}

		// Received an error, didn't match
		if err, ok := v.(error); ok {
			return false, err
		}

		// Didn't match, skip
		if v == false {
			return false, nil
		}
	}

	return true, nil
}

func routeWebhooks(body []byte, event string, routes []route) error {
	for _, r := range routes {
		matches, err := routeMatches(body, event, r)
		if err != nil {
			fmt.Printf("%s had an error, skipping\n", err)
		}

		if matches {
			routeWebhook(body, r)
		}
	}

	return nil
}

func handleWebhook(body []byte, signature string, secret string, event string, routes []route) (int, error) {
	validSignature, err := isValidSignature(secret, signature, body)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to check if signature is valid: %w+", err)
	}

	if !validSignature {
		return http.StatusBadRequest, fmt.Errorf("invalid signature")
	}

	err = routeWebhooks(body, event, routes)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to route webhook: %w+", err)
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

	routes := getRoutes(os.Environ())

	statusCode, err := handleWebhook(bs, signature, secret, r.Header.Get("X-GitHub-Event"), routes)
	if err != nil {
		errorResponse(w, r, err, statusCode)
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": "Success"})
}
