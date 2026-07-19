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
	NoLoginErr      = NewAuthError("User is not logged in.", errors.New("user is not logged in"))
	UnAuthorizedErr = NewAuthError("User does not have access to this resource.",
		errors.New("user is not authorized"))

	UsernameInUseError = NewDBError(http.StatusBadRequest, "Username is already in use",
		fmt.Errorf("DB Error: %v", pgerrcode.UniqueViolation))
	UserNotFoundError = NewDBError(http.StatusBadRequest, "Username is already in use",
		pgx.ErrNoRows)
)
