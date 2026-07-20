package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int

	HttpServer *http.Server
}

const HomeURL = "https://localhost:8080"

func NewServer() *Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	NewServer := &Server{
		port: port,
	}

	// Declare Server config
	NewServer.HttpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return NewServer
}
