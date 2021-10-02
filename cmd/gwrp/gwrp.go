package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/shteou/gwrp/pkg/handlers"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/status", handlers.StatusHandler)
	r.HandleFunc("/webhook", handlers.WebhookHandler).
		Methods("POST")

	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}
	srv.ListenAndServe()
}
