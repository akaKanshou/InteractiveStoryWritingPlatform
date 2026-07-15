package server

import (
	"crypto/rand"

	"github.com/gin-gonic/gin"
)

type story struct {
	StoryID     string `json:"story_id" db:"story_id"`
	StoryName   string `json:"story_name" db:"story_name"`
	Description string `json:"description" db:"description"`
}

func (s *Server) newStory(c *gin.Context) error {
	st := story{
		StoryID: rand.Text()[:15],
	}

	st.Description = ""

	return nil
}
