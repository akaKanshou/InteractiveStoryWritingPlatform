package server

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Story struct {
	StoryID     string `json:"story_id" db:"story_id"`
	StoryName   string `json:"story_name" db:"story_name"`
	Description string `json:"description" db:"description"`
	Visibility  int    `json:"visibility" db:"visibility"`
	Username    string `json:"username" db:"username"`
}

type Chapter struct {
	ChapterID   string `json:"chapter_id" db:"chapter_id"`
	ChapterName string `json:"chapter_name" db:"chapter_name"`
	FileID      string `json:"-" db:"file_id"`
	StoryID     string `json:"story_id" db:"story_id"`
	Content     string `json:"content,omitempty" db:"content"`
}

const (
	visibilityPrivate int8 = iota
	visibilityPublic
)

func validateStory(story *Story) error {
	if len(story.StoryName) < 1 {
		return errors.New("story name can not be empty")
	}

	if len(story.StoryName) > 50 {
		return errors.New("story name can not be longer than 50 characters")
	}

	if (story.Visibility < 0) || (story.Visibility > 1) {
		return errors.New("invalid story visibility")
	}

	return nil
}

func (s *Server) newStory(c *gin.Context, story *Story) (string, error) {
	if err := validateStory(story); err != nil {
		return "", err
	}

	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		return "", err
	}

	if err = s.insertStory(story, &u); err != nil {
		return "", err
	}

	return story.StoryID, nil
}

func (s *Server) getStories(username, visibility string, c *gin.Context) ([]Story, error) {
	if validateUsername(username) != nil {
		return nil, fmt.Errorf("invalid username")
	}

	if visibility == "public" {
		return s.getAllUserStoriesFromDBWithVisibility(username, visibilityPublic)

	}

	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		return nil, err
	}

	if (visibility != "private") && (visibility != "all") {
		return nil, fmt.Errorf("invalid visibility")
	}

	if !(u.getState(AuthGoogle) && u.getState(AuthDB) && (u.Username == username)) {
		return nil, errors.New("unauthorized access")
	}

	if visibility == "private" {
		return s.getAllUserStoriesFromDBWithVisibility(username, visibilityPrivate)
	}

	return s.getAllUserStoriesFromDB(username)
}

func (s *Server) editStory(c *gin.Context, newSt *Story) error {
	oldSt, err := s.getStory(newSt.StoryID)

	if err != nil {
		return err
	}

	if newSt.Description != "" {
		oldSt.Description = newSt.Description
	}

	if newSt.StoryName != "" {
		oldSt.StoryName = newSt.StoryName
	}

	if newSt.Visibility != -1 {
		oldSt.Visibility = newSt.Visibility
	}

	if err := validateStory(oldSt); err != nil {
		return err
	}

	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		return err
	}

	err = s.updateStory(oldSt, &u)
	return err
}

func validateChapter(chapter *Chapter) error {
	if len(chapter.ChapterName) > 50 {
		return errors.New("chapter name can not be longer than 50 characters")
	}

	if chapter.StoryID == "" {
		return errors.New("invalid chapter story id")
	}

	return nil
}

func (s *Server) newChapter(c *gin.Context, chapter *Chapter, content string) (string, error) {
	if err := validateChapter(chapter); err != nil {
		return "", err
	}

	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		return "", err
	}

	if err = s.insertChapter(chapter, &u); err != nil {
		return "", err
	}

	if err = saveChapterContentToFile(chapter, content); err != nil {
		return "", err
	}

	return chapter.ChapterID, nil
}

func (s *Server) editChapter(c *gin.Context, newChapter *Chapter, content string) error {
	chapter, err := s.getChapter(newChapter.ChapterID)

	if err != nil {
		return err
	}

	if newChapter.ChapterName != "" {
		chapter.ChapterName = newChapter.ChapterName
	}

	if err := validateChapter(chapter); err != nil {
		return err
	}

	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		return err
	}

	if err = s.updateChapter(chapter, &u); err != nil {
		return err
	}

	if content != "" {
		err = saveChapterContentToFile(chapter, content)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) getStoryChapters(storyID string, c *gin.Context) ([]Chapter, error) {
	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		return nil, err
	}

	story, err := s.getStory(storyID)
	if err != nil {
		return nil, err
	}

	if (story.Visibility == 0) && (story.Username != u.Username) {
		return nil, errors.New("unauthorized access")
	}

	return s.getChaptersByStory(storyID)
}

func (s *Server) getChapterWithContent(chapterID string, c *gin.Context) (*Chapter, error) {
	chapter, err := s.getChapter(chapterID)
	if err != nil {
		return nil, err
	}

	story, err := s.getStory(chapter.StoryID)
	if err != nil {
		return nil, err
	}

	u, err := s.getUserInfoFromSession(c)
	if err != nil {
		return nil, err
	}

	if (story.Visibility == 0) && (u.Username != story.Username) {
		return nil, errors.New("unauthorized access")
	}

	chapter.Content, err = getContent(chapter.FileID)
	if err != nil {
		return nil, err
	}

	return chapter, nil
}
