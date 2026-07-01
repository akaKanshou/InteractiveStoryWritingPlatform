package server

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	_ "github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Server struct {
	port int
}

var store *sessions.CookieStore
var googleOAuthConfig *oauth2.Config

type User struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	GoogleUserID string `json:"id"`
}

func init() {
	gob.Register(User{})

	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.MaxAge(86400)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
	}

	googleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENTID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENTSECRET"),
		RedirectURL:  "https://localhost:8080/auth/google/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
