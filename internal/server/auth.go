package server

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func (s *Server) setAuthSession(c *gin.Context) (string, error) {
	session, err := s.store.Get(c.Request, "sessionISWP")

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	stateString := url.QueryEscape(rand.Text())
	session.Values["state"] = stateString

	if c.Query("remember") == "true" {
		session.Options.MaxAge = 7 * 86400
	} else {
		session.Options.MaxAge = 86400
	}

	oAuthReqUrl := s.googleOAuthConfig.AuthCodeURL(stateString, oauth2.AccessTypeOffline)

	if err := session.Save(c.Request, c.Writer); err != nil {
		fmt.Println(err)
		return "", err
	}

	return oAuthReqUrl, nil
}

func (s *Server) exchangeTokenForUser(code string, c *gin.Context) (User, error) {
	token, err := s.googleOAuthConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		fmt.Println("errExchange: ", err)
		return User{}, err
	}

	authClient := s.googleOAuthConfig.Client(context.Background(), token)

	authRespFromGoogle, err := authClient.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		fmt.Println("errResp: ", err)
		return User{}, err
	}

	defer authRespFromGoogle.Body.Close()

	var v User
	err = json.NewDecoder(authRespFromGoogle.Body).Decode(&v)
	if err != nil {
		fmt.Println("errDecoding", err)
		return User{}, err
	}

	v.AuthState |= AuthGoogle
	return v, nil
}

func (s *Server) authCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := s.getUserInfoFromSession(c)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusForbidden)
			return
		}

		if u.getState(AuthGoogle) && u.getState(AuthDB) {
			c.Next()
		} else {
			c.Status(http.StatusForbidden)
			return
		}
	}
}
