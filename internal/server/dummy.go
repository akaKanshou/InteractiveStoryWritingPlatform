package server

import (
	"crypto/rand"
	"forgeverse/internal/db"
	fverrors "forgeverse/internal/errors"
	"forgeverse/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) makeDummyUser(c *gin.Context) {
	_, err := db.RegisterUser(&models.User{
		Username: c.Query("username"),
		UserID:   "D-" + rand.Text()[0:15],
		Email:    "dummy@dummy.com",
	})

	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func (s *Server) newDummyStoryHandler(c *gin.Context) {
	user := &models.User{
		Username: c.PostForm("username"),
	}

	storyID, err := db.InsertNewStory(c.PostForm("story_name"), c.PostForm("description"), rand.Text()[0:15],
		models.VisibilityPublic, user)

	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"story_id": storyID})
}

func (s *Server) newDummyChapterHandler(c *gin.Context) {
	chapter := &models.Chapter{
		ChapterID:   rand.Text()[:15],
		ChapterName: c.PostForm("chapter_name"),
		StoryID:     c.PostForm("story_id"),
		Content:     c.PostForm("content"),
		FileID:      rand.Text()[:15],
	}

	err := db.InsertNewChapter(chapter)
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"chapter": chapter})
}

func (s *Server) addDummyEdgeHandler(c *gin.Context) {
	err := db.InsertNewEdge(c.PostForm("from_chapter"), c.PostForm("to_chapter"))
	if err != nil {
		fverrors.SendErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}
