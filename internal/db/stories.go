package db

import (
	"context"
	"errors"
	"fmt"
	fverrors "forgeverse/internal/errors"
	"forgeverse/internal/models"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// GetStoriesByUser gets all stories of given user with specified visibility
func GetStoriesByUser(username string, visibilityVal models.StoryVisibility) ([]models.Story,
	fverrors.Error) {
	rows, pgErr := dbConn.Query(context.Background(),
		`SELECT story_id, story_name, 
description, visibility FROM stories WHERE username=$1 AND visibility&$2>0`,
		username, visibilityVal)

	if pgErr != nil {
		return nil, fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred.", pgErr)
	}

	defer rows.Close()

	stories := make([]models.Story, 0, 50)
	for rows.Next() {
		story := models.Story{
			Username: username,
		}

		if err := rows.Scan(&story.StoryID, &story.StoryName, &story.Description, &story.Visibility); err != nil {
			return nil, fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred.", err)
		}

		stories = append(stories, story)
	}

	if err := rows.Err(); err != nil {
		return nil, fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred.", err)
	}

	return stories, nil
}

// GetStoryByID gets story by its story_id. Returns StoryNotFoundError if no such story exists.
func GetStoryByID(storyID string) (*models.Story, fverrors.Error) {
	story := &models.Story{
		StoryID: storyID,
	}

	pgErr := dbConn.QueryRow(context.Background(),
		"SELECT username, story_name, description, visibility FROM stories WHERE story_id = $1",
		storyID).Scan(
		&story.Username,
		&story.StoryName,
		&story.Description,
		&story.Visibility,
	)
	if pgErr != nil && errors.Is(pgErr, pgx.ErrNoRows) {
		return nil, fverrors.StoryNotFoundError
	} else if pgErr != nil {
		return nil, fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred.", pgErr)
	}

	return story, nil
}

// InsertNewStory inserts a new story with given details.
func InsertNewStory(storyName, storyDescription, storyID string, visibility models.StoryVisibility,
	user *models.User) (string,
	fverrors.Error) {
	_, err := dbConn.Exec(context.Background(),
		`INSERT INTO stories (story_name, description, visibility, username, story_id) VALUES ($1, $2, $3,$4 ,$5)`,
		storyName, storyDescription, visibility, user.Username, storyID)

	if err == nil {
		return storyID, nil
	}

	if pgErr, okay := errors.AsType[*pgconn.PgError](err); okay && (pgErr.Code == pgerrcode.UniqueViolation) {
		return "", fverrors.NewInvalidRequestError(fmt.Sprintf("User already has a story named %s", storyName),
			pgErr)
	} else if okay && (pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) || pgerrcode.IsDataException(pgErr.Code)) {
		return "", fverrors.GenericInvalidRequestErr
	}

	return "", fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred.", err)
}

// UpdateStory all description, visibility and story_name for given story_id.
// New values from these fields are taken from "story" argument. story.Username is not used.
func UpdateStory(story *models.Story) fverrors.Error {
	_, err := dbConn.Exec(context.Background(),
		"UPDATE stories SET description=$1, visibility=$2, story_name=$3 WHERE story_id=$4",
		story.Description, story.Visibility, story.StoryName, story.StoryID)

	if pgErr, okay := errors.AsType[*pgconn.PgError](err); okay && (pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) || pgerrcode.IsDataException(pgErr.Code)) {
		return fverrors.GenericInvalidRequestErr
	} else if err != nil {
		return fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred.", err)
	}

	return nil
}
