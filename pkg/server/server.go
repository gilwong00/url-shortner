package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

type Server struct {
	Port  string
	Store *redis.Client
}

func NewServer(port string, store *redis.Client) *Server {
	return &Server{
		Port:  port,
		Store: store,
	}
}

func (s *Server) StartServer() {
	mux := http.NewServeMux()
	// add routes
	server := http.Server{
		Addr:         fmt.Sprintf(":%s", s.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	// start the server
	go func() {
		fmt.Printf("Starting server on port %s\n", s.Port)
		err := server.ListenAndServe()
		if err != nil {
			fmt.Printf("Error starting server: %s", err.Error())
			os.Exit(1)
		}
	}()
	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)
	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
