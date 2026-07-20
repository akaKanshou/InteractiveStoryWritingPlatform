package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
)

type Error interface {
	ResponseError() string
	Code() int
	Error() string
}

type jsonError struct {
	Message string `json:"message"`
}

func SendErrorResponse(c *gin.Context, err Error) {
	c.JSON(err.Code(), jsonError{err.ResponseError()})
	fmt.Println(err.Error())
}

var (
	NoLoginErr      = NewAuthError("User is not logged in", errors.New("user is not logged in"))
	UnAuthorizedErr = NewAuthError("User does not have access to this resource",
		errors.New("user is not authorized"))

	GenericInvalidRequestErr = NewAuthError("Invalid request parameters", errors.New("invalid request"))

	UsernameInUseError = NewDBError(http.StatusUnprocessableEntity, "Username is already in use",
		fmt.Errorf("DB Error: %v", pgerrcode.UniqueViolation))
	UserNotFoundError = NewDBError(http.StatusNotFound, "No such user exists",
		pgx.ErrNoRows)
	StoryNotFoundError   = NewDBError(http.StatusNotFound, "No such story exists", pgx.ErrNoRows)
	ChapterNotFoundError = NewDBError(http.StatusNotFound, "No such chapter exists", pgx.ErrNoRows)
)
