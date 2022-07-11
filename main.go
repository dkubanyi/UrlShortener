package main

import (
	"context"
	"dkubanyi/urlShortener/handler"
	"dkubanyi/urlShortener/storage/redis"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	serverHostName  = "SERVER_HOST"
	serverPortName  = "SERVER_PORT"
	accessTokenName = "ACCESS_TOKEN"
	redisHostName   = "REDIS_HOST"
	redisPortName   = "REDIS_PORT"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("Environment variables failed to load")
	}

	host := os.Getenv(serverHostName)
	port := os.Getenv(serverPortName)
	accessToken := os.Getenv(accessTokenName)
	redisHost := os.Getenv(redisHostName)
	redisPort := os.Getenv(redisPortName)

	redis, err := redisStorage.New(redisHost, redisPort)
	if err != nil {
		log.Fatalf("%v", err)
	}

	defer redis.Close()

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: handler.New(host+":"+port, accessToken, redis),
	}

	go func() {
		log.Printf("Starting HTTP server at %q", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("%v", err)
		} else {
			log.Println("Server closed!")
		}
	}()

	// Check for a closing signal & graceful shutdown
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Interrupt, syscall.SIGTERM)

	sig := <-sigquit
	log.Printf("Caught sig: %+v", sig)
	log.Printf("Gracefully shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("Unable to shut down server: %v", err)
	} else {
		log.Println("Server stopped")
	}
}
