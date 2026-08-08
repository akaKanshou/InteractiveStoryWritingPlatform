package api

import (
	"crypto/rand"
	"errors"
	"forgeverse/internal/auth"
	"forgeverse/internal/db"
	fverrors "forgeverse/internal/errors"
	"forgeverse/internal/models"

	"github.com/gin-gonic/gin"
)

func GetStoriesByUser(user *models.User, c *gin.Context, username, visibility string) ([]models.Story, fverrors.Error) {
	if err := models.ValidateUsername(username); err != nil {
		return nil, err
	}

	target, err := db.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	visibilityVal, err := models.ValidateNConvertVisibility(visibility)
	if err != nil {
		return nil, err
	}

	if (visibilityVal&models.VisibilityPrivate > 0) && (!auth.IsPrivateAuthenticated(user) || (user.Username != target.Username)) {
		return nil, fverrors.UnAuthorizedErr
	}

	return db.GetStoriesByUser(username, visibilityVal)
}

func CreateNewStory(c *gin.Context, user *models.User) (string, fverrors.Error) {
	if !auth.IsPrivateAuthenticated(user) {
		return "", fverrors.UnAuthorizedErr
	}

	storyName, visibility, description := c.PostForm("story_name"), c.PostForm("visibility"), c.PostForm("description")

	if storyName == "" {
		return "", fverrors.NewInvalidRequestError("Story name can not be empty", errors.New("story name can not be empty"))
	}

	visibilityVal, err := models.ValidateNConvertVisibility(visibility)
	if err != nil {
		return "", err
	}

	return db.InsertNewStory(storyName, description, rand.Text()[:15], visibilityVal, c.PostForm("forkable") == "true", user)
}

func EditStory(c *gin.Context, user *models.User) fverrors.Error {
	storyID := c.PostForm("story_id")
	if err := models.ValidateRID(storyID); err != nil {
		return err
	}

	story, err := db.GetStoryByID(storyID)
	if err != nil {
		return err
	}

	if story.Username != user.Username {
		return fverrors.UnAuthorizedErr
	}

	story.Visibility, err = models.ValidateNConvertVisibility(c.PostForm("visibility"))
	if err != nil {
		return err
	}

	story.StoryName, story.Description = c.PostForm("story_name"), c.PostForm("description")
	story.Forkable = c.PostForm("forkable") == "true"

	if err = db.UpdateStory(story); err != nil {
		return err
	}

	return nil
}

func GetHomePageStories() ([][]models.Story, fverrors.Error) {
	stories := make([][]models.Story, 3)
	var err fverrors.Error

	stories[0], err = db.GetLatestStories(3)
	if err != nil {
		return nil, err
	}

	stories[1] = nil

	stories[2], err = db.GetLatestStories(21)
	if err != nil {
		return nil, err
	}

	return stories, nil
}
