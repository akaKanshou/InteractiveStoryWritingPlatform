package server

import (
	"crypto/rand"
	"errors"
	"unicode"

	"github.com/gin-gonic/gin"
)

type story struct {
	StoryID     string `json:"story_id" db:"story_id"`
	StoryName   string `json:"story_name" db:"story_name"`
	Description string `json:"description" db:"description"`
}

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

	return nil
}

func (s *Server) newStory(c *gin.Context) (string, error) {
	st := story{
		StoryID:     rand.Text()[:15],
		StoryName:   c.PostForm("story_name"),
		Description: c.PostForm("description"),
	}

	if err := validateStory(&st); err != nil {
		return "", err
	}

	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		return "", err
	}

	if err = s.insertStory(&st, &u); err != nil {
		return "", err
	}

	return st.StoryID, nil
}
