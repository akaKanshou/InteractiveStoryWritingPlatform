package models

import (
	"errors"
	"fmt"
	fverrors "forgeverse/internal/errors"
	"regexp"
)

var (
	usernameRegex  *regexp.Regexp
	storyIDRegex   *regexp.Regexp
	chapterIDRegex *regexp.Regexp
)

func init() {
	var err error

	usernameRegex, err = regexp.Compile("^[a-z0-9_]{3,18}$")
	if err != nil {
		panic(err)
	}

	storyIDRegex, err = regexp.Compile("^[A-Z0-9]{15}$")
	if err != nil {
		panic(err)
	}

	chapterIDRegex, err = regexp.Compile("^[A-Z0-9]{15}$")
	if err != nil {
		panic(err)
	}
}

func ValidateUsername(username string) fverrors.Error {
	if !usernameRegex.MatchString(username) {
		return fverrors.NewInvalidRequestError("Invalid username.", errors.New("failed regex: models.usernameRegex"))
	}

	return nil
}

func ValidateStoryID(storyID string) fverrors.Error {
	if !storyIDRegex.MatchString(storyID) {
		return fverrors.NewInvalidRequestError("Invalid username.", errors.New("failed regex: models.storyIDRegex"))
	}

	return nil
}

func ValidateChapterID(chapterID string) fverrors.Error {
	if !chapterIDRegex.MatchString(chapterID) {
		return fverrors.NewInvalidRequestError("Invalid chapter ID.", errors.New("failed regex: models.chapterIDRegex"))
	}

	return nil
}

func ValidateVisibility(visibility string) (StoryVisibility, fverrors.Error) {
	switch visibility {
	case "private":
		return VisibilityPrivate, nil
	case "public":
		return VisibilityPublic, nil
	case "all":
		return VisibilityPrivate | VisibilityPublic, nil
	}

	return 0, fverrors.NewInvalidRequestError("Invalid visibility.", fmt.Errorf("invalid visibility \"%s\"", visibility))
}
