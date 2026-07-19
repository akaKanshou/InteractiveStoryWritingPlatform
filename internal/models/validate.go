package models

import (
	"errors"
	"fmt"
	fverrors "forgeverse/internal/errors"
	"regexp"
)

var (
	usernameRegex *regexp.Regexp
)

func init() {
	var err error

	usernameRegex, err = regexp.Compile("^[a-z0-9_]{3-18}$")
	if err != nil {
		panic(err)
	}

}

func ValidateUsername(username string) fverrors.Error {
	if !usernameRegex.MatchString(username) {
		return fverrors.NewBadRequestError("Invalid username.", errors.New("failed regex: models.usernameRegex"))
	}

	return nil
}

func ValidateVisibility(visibility string) (int8, fverrors.Error) {
	switch visibility {
	case "private":
		return visibilityPrivate, nil
	case "public":
		return visibilityPublic, nil
	case "all":
		return visibilityPrivate | visibilityPublic, nil
	}

	return 0, fverrors.NewBadRequestError("Invalid visibility.", fmt.Errorf("invalid visibility \"%s\"", visibility))
}
