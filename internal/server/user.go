package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	AuthGoogle uint8 = 0x1
	AuthDB     uint8 = 0x2
	AuthAny    uint8 = 0xFF
)

type user struct {
	Name         string `json:"name" db:"display_name"`
	Email        string `json:"email" db:"email"`
	GoogleUserID string `json:"id" db:"user_id"`
	Username     string `db:"username"`

	AuthState uint8
}

func (s *Server) getUserInfoFromSession(c *gin.Context) (user, error) {
	session, err := s.store.Get(c.Request, "sessionISWP")
	if err != nil {
		fmt.Println(err)
		return user{}, err
	}

	userInfo := session.Values["userInfo"]

	if userInfo == nil {
		return user{Name: "Stranger"}, nil
	}

	return userInfo.(user), nil
}

func (s *Server) setUserInfoToSession(c *gin.Context, u user) error {
	session, err := s.store.Get(c.Request, "sessionISWP")
	if err != nil {
		fmt.Println(err)
		return err
	}

	session.Values["userInfo"] = u
	if err = session.Save(c.Request, c.Writer); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (v user) getState(state uint8) bool {
	return v.AuthState&state > 0
}
