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

func GetStoryByID(user *models.User, c *gin.Context) (*models.Story, fverrors.Error) {
	story, err := db.GetStoryByID(c.Param("story_id"))
	if err != nil {
		return nil, err
	}

	if (story.Visibility == models.VisibilityPrivate) && (user.Username != story.Username) {
		return nil, fverrors.UnAuthorizedErr
	}

	return story, nil
}

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

func GetHomePageStories() ([]models.Story, []models.Story, fverrors.Error) {

	latestStories, err := db.GetLatestStories(4)
	if err != nil {
		return nil, nil, err
	}

	editorsPicks := make([]models.Story, 0, 4)
	for _, id := range models.EditorsPickIDs {
		story, err := db.GetStoryByID(id)
		if err != nil {
			return nil, nil, err
		}

		editorsPicks = append(editorsPicks, *story)
	}

	return nil, latestStories, err
}

func DeleteStory(user *models.User, c *gin.Context) fverrors.Error {
	storyId := c.Param("story_id")
	story, err := db.GetStoryByID(storyId)
	if err != nil {
		return err
	}

	if story.Username != user.Username {
		return fverrors.UnAuthorizedErr
	}

	return db.DeleteStory(storyId)
}
