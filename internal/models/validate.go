package models

import (
	"errors"
	"fmt"
	fverrors "forgeverse/internal/errors"
	"regexp"
	"strconv"
)

var (
	usernameRegex  *regexp.Regexp
	rIDRegex       *regexp.Regexp
	chapterIDRegex *regexp.Regexp
)

func init() {
	var err error

	usernameRegex, err = regexp.Compile("^[a-z0-9_]{3,18}$")
	if err != nil {
		panic(err)
	}

	rIDRegex, err = regexp.Compile("^[A-Z0-9]{15}$")
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

func ValidateRID(rID string) fverrors.Error {
	if !rIDRegex.MatchString(rID) {
		return fverrors.NewInvalidRequestError("Invalid ID.", errors.New("failed regex: models.rIDRegex"))
	}

	return nil
}

func ValidateNConvertVisibility(visibility string) (StoryVisibility, fverrors.Error) {
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

func ValidateNConvertPage(page string) (int, fverrors.Error) {
	if page == "" {
		return 0, nil
	}

	n, err := strconv.Atoi(page)
	if err != nil {
		return 0, fverrors.NewInvalidRequestError("Invalid page.", fmt.Errorf("invalid page \"%s\"", page))
	}

	return n, nil
}
