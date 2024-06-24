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

	"github.com/gilwong00/url-shortner/pkg/config"
	"github.com/gilwong00/url-shortner/pkg/handlers"
	"github.com/gilwong00/url-shortner/pkg/middleware"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	Port  int
	Store *redis.Client
}

func NewServer(port int, store *redis.Client) *Server {
	return &Server{
		Port:  port,
		Store: store,
	}
}

func (s *Server) StartServer(config *config.Config) {
	handler := handlers.NewHandler(config, s.Store)
	mux := http.NewServeMux()
	// routes
	// GET
	mux.HandleFunc("GET /url/{shortName}", handler.GetURL)
	// POST
	mux.Handle("POST /url", middleware.RateLimiter(
		http.HandlerFunc(handler.CreateShortenURL),
		context.Background(),
		s.Store,
		config.MaxRequestLimit,
	),
	)
	// PUT
	mux.HandleFunc("PUT /url/{shortName}", handler.UpdateURL)
	// DELETE
	mux.HandleFunc("DELETE /url/{shortName}", handler.DeleteURL)

	server := http.Server{
		Addr:         fmt.Sprintf(":%v", s.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	// start the server
	go func() {
		fmt.Printf("Starting server on port %v\n", s.Port)
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
