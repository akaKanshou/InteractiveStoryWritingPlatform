package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name         string `json:"name" db:"display_name"`
	Email        string `json:"email" db:"email"`
	GoogleUserID string `json:"id" db:"user_id"`

	Username string `db:"username"`
}

func (s *Server) getUserInfoFromSession(c *gin.Context) (User, error) {
	session, err := s.store.Get(c.Request, "sessionISWP")
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	userInfo := session.Values["userInfo"]

	if userInfo == nil {
		return User{Name: "Stranger"}, nil
	}

	return userInfo.(User), nil
}

func (s *Server) setUserInfoToSession(c *gin.Context, u User) error {
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
