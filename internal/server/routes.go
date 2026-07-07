package server

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost:8080"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/", s.HelloWorldHandler)
	r.GET("/auth/google", s.authGoogleHandler)
	r.GET("/auth/google/callback", s.authGoogleCallbackHandler)
	r.GET("/logout", s.logoutHandler)

	return r
}

func (s *Server) getUserInfoFromSession(c *gin.Context) (*User, error) {
	session, err := s.store.Get(c.Request, "sessionISWP")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	userInfo := session.Values["userInfo"]

	if userInfo == nil {
		return nil, nil
	}

	return userInfo.(*User), nil
}

func (s *Server) logoutHandler(c *gin.Context) {
	session, errGettingState := s.store.Get(c.Request, "sessionISWP")

	if errGettingState != nil {
		fmt.Println(errGettingState)
		c.Writer.WriteHeader(500)
		return
	}

	session.Options.MaxAge = -1
	session.Values = make(map[any]any)
	errSaving := session.Save(c.Request, c.Writer)
	if errSaving != nil {
		fmt.Println(errSaving)
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (s *Server) authGoogleCallbackHandler(c *gin.Context) {
	code, stateString := c.Query("code"), c.Query("state")

	if code == "" {
		fmt.Println("Invalid callback")
		c.Writer.WriteHeader(400)
		return
	}

	session, errGettingState := s.store.Get(c.Request, "sessionISWP")

	if errGettingState != nil {
		fmt.Println(errGettingState)
		c.Writer.WriteHeader(500)
		return
	}

	if session.Values["state"] != stateString {
		c.Writer.WriteHeader(401)
		return
	}

	token, errExchange := s.googleOAuthConfig.Exchange(c.Request.Context(), code)
	if errExchange != nil {
		fmt.Println("errExchange: ", errExchange)
		http.Error(c.Writer, errExchange.Error(), http.StatusBadRequest)
		return
	}

	authClient := s.googleOAuthConfig.Client(context.Background(), token)

	authRespFromGoogle, errResp := authClient.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if errResp != nil {
		fmt.Println("errResp: ", errResp)
		c.Writer.WriteHeader(500)
		return
	}

	defer authRespFromGoogle.Body.Close()

	var v User
	errDecoding := json.NewDecoder(authRespFromGoogle.Body).Decode(&v)
	if errDecoding != nil {
		fmt.Println("errDecoding", errDecoding)
		c.Writer.WriteHeader(500)
		return
	}

	session.Values["userInfo"] = v

	if errSaving := session.Save(c.Request, c.Writer); errSaving != nil {
		fmt.Println(errSaving)
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (s *Server) authGoogleHandler(c *gin.Context) {
	userInfo, err := s.getUserInfoFromSession(c)

	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(500)
		return
	} else if userInfo != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	session, err := s.store.Get(c.Request, "sessionISWP")

	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(500)
		return
	}

	stateString := url.QueryEscape(rand.Text())
	session.Values["state"] = stateString

	oAuthReqUrl := s.googleOAuthConfig.AuthCodeURL(stateString, oauth2.AccessTypeOffline)

	if errSaving := session.Save(c.Request, c.Writer); errSaving != nil {
		fmt.Println(errSaving)
	}

	http.Redirect(c.Writer, c.Request, oAuthReqUrl, http.StatusTemporaryRedirect)
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	t, err := template.ParseFiles("./internal/webpages/index.html")

	userInfo, err := s.getUserInfoFromSession(c)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(500)
		return
	}

	if userInfo == nil {
		userInfo = &User{
			Name: "Stranger",
		}
	}

	if err = t.Execute(c.Writer, *userInfo); err != nil {
		fmt.Println(err)
	}
}
