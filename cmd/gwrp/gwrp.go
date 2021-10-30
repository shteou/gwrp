package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	gorillahandlers "github.com/gorilla/handlers"
	"github.com/shteou/gwrp/pkg/handlers"
)

type server struct {
	http.Server
	reqCount uint32
}

func (s *server) waitShutdown(channel chan os.Signal) {
	// listen for SiGINT (ctrl+c) and SIGTERM
	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)

	sig := <-channel
	fmt.Printf("Received interrupt (signal: %v)\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := s.Shutdown(ctx)
	if err != nil {
		fmt.Printf("Failed to shutdown: %v\n", err)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/status", handlers.StatusHandler)
	r.HandleFunc("/webhook", handlers.WebhookHandler).
		Methods("POST")

	http.Handle("/", r)

	loggingHandler := gorillahandlers.LoggingHandler(os.Stdout, r)

	srv := server{
		Server: http.Server{
			Handler:      loggingHandler,
			Addr:         "0.0.0.0:8080",
			WriteTimeout: 60 * time.Second,
			ReadTimeout:  60 * time.Second,
		},
	}

	done := make(chan os.Signal, 1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Printf("Listen and serve: %v\n", err)
		}
		// Emulate shutdown interrupt
		done <- syscall.SIGINT
	}()

	srv.waitShutdown(done)
}
