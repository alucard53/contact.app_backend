package main

import (
	"log"
	"net/http"

	"context"
	"os"
	"os/signal"
	"time"

	"contact.app_backend/handlers"
	"github.com/fjl/go-couchdb"
	"github.com/rs/cors"
)

type Server struct {
	srv *http.Server
	l   *log.Logger
}

func initServer(l *log.Logger, db *couchdb.DB) Server {

	mx := http.NewServeMux()

	mx.Handle("/contacts", handlers.NewContact(l, db))

	return Server{
		srv: &http.Server{
			Addr:    ":6969",
			Handler: cors.AllowAll().Handler(mx),
		},
		l: l,
	}
}

func startServer(server Server) {
	go func() {
		err := server.srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	handleShutdown(server)
}

func handleShutdown(server Server) {
	sigch := make(chan os.Signal, 1)

	signal.Notify(sigch, os.Interrupt)

	sig := <-sigch

	server.l.Println("Got signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	server.srv.Shutdown(ctx)
}
