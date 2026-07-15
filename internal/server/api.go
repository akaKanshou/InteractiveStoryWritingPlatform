package server

import (
	"errors"
	"unicode"

	"github.com/gin-gonic/gin"
)

type story struct {
	StoryID     string `json:"story_id" db:"story_id"`
	StoryName   string `json:"story_name" db:"story_name"`
	Description string `json:"description" db:"description"`
	Visibility  int    `json:"visibility" db:"visibility"`
}

const (
	visibiltyPrivate int8 = iota
	visibiltyPublic
)

func validateStory(st *story) error {
	if len(st.StoryName) < 1 {
		return errors.New("story name can not be empty")
	}

	for _, c := range st.StoryName {
		if c > unicode.MaxASCII {
			return errors.New("story name contains invalid characters")
		}
	}

	if len(st.StoryName) > 50 {
		return errors.New("story name can not be longer than 50 characters")
	}

	if (st.Visibility < 0) || (st.Visibility > 1) {
		return errors.New("invalid story visibility")
	}

	return nil
}

func (s *Server) newStory(c *gin.Context, st *story) (string, error) {
	if err := validateStory(st); err != nil {
		return "", err
	}

	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		return "", err
	}

	if err = s.insertStory(st, &u); err != nil {
		return "", err
	}

	return st.StoryID, nil
}

func (s *Server) editStory(c *gin.Context, st *story) error {
	if err := validateStory(st); err != nil {
		return err
	}

	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		return err
	}

	if err = s.updateStory(st, &u); err != nil {
		return err
	}

	return nil
}
