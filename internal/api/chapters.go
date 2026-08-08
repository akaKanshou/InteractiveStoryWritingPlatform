package api

import (
	"crypto/rand"
	"forgeverse/internal/db"
	fverrors "forgeverse/internal/errors"
	"forgeverse/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateNewChapter(c *gin.Context, user *models.User) (string, fverrors.Error) {
	story, err := db.GetStoryByID(c.PostForm("story_id"))
	if err != nil {
		return "", err
	}

	chapter := &models.Chapter{
		ChapterID:   rand.Text()[:15],
		ChapterName: c.PostForm("chapter_name"),
		StoryID:     story.StoryID,
		Content:     c.PostForm("content"),
		FileID:      rand.Text()[:15],
		Username:    c.PostForm("username"),
	}

	timeNow := db.TimeNow()

	err = db.InsertNewChapter(chapter, timeNow)
	if err != nil {
		return "", err
	}

	if err = db.WriteToFile(chapter.FileID, chapter.Content); err != nil {
		return "", err
	}

	return chapter.ChapterID, nil
}

func EditChapter(c *gin.Context, user *models.User) fverrors.Error {
	chapterID := c.PostForm("chapter_id")
	if err := models.ValidateRID(chapterID); err != nil {
		return fverrors.GenericInvalidRequestErr
	}

	chapter, err := db.GetChapterByID(chapterID)
	if err != nil {
		return err
	}

	story, err := db.GetStoryByID(chapter.StoryID)
	if err != nil {
		return err
	}

	if story.Username != user.Username {
		return fverrors.UnAuthorizedErr
	}

	timeNow := db.TimeNow()

	chapter.ChapterName = c.PostForm("chapter_name")
	if err := db.UpdateChapter(chapter, timeNow); err != nil {
		return err
	}

	if err := db.UpdateStoryLastUpdated(story.StoryID, timeNow); err != nil {
		return err
	}

	if err := db.WriteToFile(chapter.FileID, c.PostForm("content")); err != nil {
		return err
	}

	return nil
}

func GetChaptersByStory(user *models.User, c *gin.Context) ([]*models.Chapter, fverrors.Error) {
	storyID := c.Param("story_id")
	if err := models.ValidateRID(storyID); err != nil {
		return nil, fverrors.GenericInvalidRequestErr
	}

	story, err := db.GetStoryByID(storyID)
	if err != nil {
		return nil, err
	}

	if (story.Visibility == models.VisibilityPrivate) && (story.Username != user.Username) {
		return nil, fverrors.UnAuthorizedErr
	}

	return db.GetChaptersByStory(storyID, 0)
}

func GetChapterByID(user *models.User, c *gin.Context) (*models.Chapter, fverrors.Error) {
	chapterID := c.Param("chapter_id")
	if err := models.ValidateRID(chapterID); err != nil {
		return nil, fverrors.GenericInvalidRequestErr
	}

	chapter, err := db.GetChapterByID(chapterID)
	if err != nil {
		return nil, err
	}

	story, err := db.GetStoryByID(chapter.StoryID)
	if err != nil {
		return nil, err
	}

	if (story.Visibility == models.VisibilityPrivate) && (story.Username != user.Username) {
		return nil, fverrors.UnAuthorizedErr
	}

	chapter.Content, err = db.ReadFromFile(chapter.FileID)

	return chapter, err
}

func CreateNewEdge(c *gin.Context, user *models.User) fverrors.Error {
	from, to := c.PostForm("from_chapter"), c.PostForm("to_chapter")
	if err := models.ValidateRID(from); err != nil {
		return fverrors.GenericInvalidRequestErr
	}
	if err := models.ValidateRID(to); err != nil {
		return fverrors.GenericInvalidRequestErr
	}

	fromChapter, err := db.GetChapterByID(from)
	if err != nil {
		return err
	}

	story, err := db.GetStoryByID(fromChapter.StoryID)
	if err != nil {
		return err
	}

	if (story.Username != user.Username) && (!story.Forkable) {
		return fverrors.UnAuthorizedErr
	}

	if err := db.InsertNewEdge(from, to); err != nil {
		return err
	}

	return nil
}

func GetEdgesByChapter(user *models.User, c *gin.Context) ([]*models.Edge, fverrors.Error) {
	chapterID := c.Param("chapter_id")
	if err := models.ValidateRID(chapterID); err != nil {
		return nil, fverrors.GenericInvalidRequestErr
	}

	chapter, err := db.GetChapterByID(chapterID)
	if err != nil {
		return nil, err
	}

	story, err := db.GetStoryByID(chapter.StoryID)
	if err != nil {
		return nil, err
	}

	if (story.Visibility == models.VisibilityPrivate) && (story.Username != user.Username) {
		return nil, fverrors.UnAuthorizedErr
	}

	return db.GetEdgesByChapter(chapterID)
}
