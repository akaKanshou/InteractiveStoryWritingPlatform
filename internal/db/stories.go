package db

import (
	"context"
	fverrors "forgeverse/internal/errors"
	"forgeverse/internal/models"
	"net/http"
)

func GetStoriesByUser(user *models.User, username, visibility string) ([]models.Story, fverrors.Error) {
	if err := models.ValidateUsername(username); err != nil {
		return nil, err
	}

	visibilityVal, err := models.ValidateVisibility(visibility)
	if err != nil {
		return nil, err
	}

	rows, pgErr := dbConn.Query(context.Background(),
		`SELECT story_id, story_name, description FROM stories WHERE username=$1 AND visibility=$2 LIMIT 50 OFFSET 0`,
		username, visibilityVal)

	if pgErr != nil {
		return nil, fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred.", err)
	}

	defer rows.Close()

	stories := make([]models.Story, 0, 50)
	for rows.Next() {
		story := models.Story{
			Username: username,
		}

		if err := rows.Scan(&story.StoryID, &story.StoryName, &story.Description); err != nil {
			return nil, fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred.", err)
		}

		stories = append(stories, story)
	}

	if err := rows.Err(); err != nil {
		return nil, fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred.", err)
	}

	return stories, nil
}

func InsertNewStory(story *models.Story) (string, error) {

	return story.StoryID, nil
}
