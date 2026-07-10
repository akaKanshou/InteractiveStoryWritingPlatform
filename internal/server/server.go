package server

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Server struct {
	port int

	HttpServer *http.Server

	store *sessions.CookieStore

	googleOAuthConfig *oauth2.Config

	PgConn *pgx.Conn
}

type User struct {
	Name         string `json:"name" db:"displayName"`
	Email        string `json:"email" db:"email"`
	GoogleUserID string `json:"id" db:"userID"`

	Username string `db:"userName"`
}

func init() {
	gob.Register(User{})
}

func NewServer() *Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	conn, err := Connect()
	if err != nil {
		log.Fatal(err)
	}

	NewServer := &Server{
		port:   port,
		store:  sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET"))),
		PgConn: conn,
	}

	NewServer.googleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENTID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENTSECRET"),
		RedirectURL:  "https://localhost:8080/auth/google/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	NewServer.store.MaxAge(7 * 86400)

	NewServer.store.Options.Path = "/"
	NewServer.store.Options.HttpOnly = true
	NewServer.store.Options.Secure = false

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
