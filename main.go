package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rajeshpudota/crypto-server/handlers"
	cache "github.com/rajeshpudota/crypto-server/internal/pkg/cache"
	"github.com/rajeshpudota/crypto-server/internal/pkg/config"
)

var (
	bindAddress string = "%s:%s"
	symbols            = []string{
		"BTCUSDT", "ETHBTC",
	}
)

const (
	Host             = "localhost"
	Port             = "8080"
	APIBaseURL       = "https://api.hitbtc.com/api/3"
	WebsocketBaseUrl = "wss://api.hitbtc.com/api/3/ws/public"
)

func main() {
	l := log.New(os.Stdout, "crypto-server", log.LstdFlags)

	config := config.Config{
		APIBaseURL:       APIBaseURL,
		WebsocketBaseUrl: WebsocketBaseUrl,
	}

	// Initialise inmemory cache
	cache := cache.NewCache(l)
	cache.UpdateCache(config, symbols)

	// Initialise validation package

	// create the handlers
	cs := handlers.NewCurrency(l, cache)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// handlers for API
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/currency/all", cs.ListAll)
	getR.HandleFunc("/currency/{symbol}", cs.ListSingle)

	//create a new http server
	host := os.Getenv("HOST")
	if host == "" {
		host = Host
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = Port
	}

	s := http.Server{
		Addr:         fmt.Sprintf(bindAddress, host, port),
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// start the server
	go func() {
		l.Printf("starting server on port: %s", port)
		if err := s.ListenAndServe(); err != nil {
			l.Printf("error starting server: %s\n", err)
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
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
