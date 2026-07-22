package db

import (
	"context"
	"errors"
	"forgeverse/internal/auth"
	fverrors "forgeverse/internal/errors"
	"forgeverse/internal/models"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// GetUserByID gets the username registered with the given userID. The returned user has no access (auth.NoAccess)
// Returns UserNotFoundError if no such user exists.
func GetUserByID(userID string) (*models.User, fverrors.Error) {
	user := &models.User{
		UserID:    userID,
		AuthState: auth.NoAccess,
	}

	err := dbConn.QueryRow(context.Background(), `SELECT username FROM users WHERE user_id=$1`,
		userID).Scan(&user.Username)

	if (err != nil) && (errors.Is(err, pgx.ErrNoRows)) {
		return nil, fverrors.UserNotFoundError
	} else if err != nil {
		return nil, fverrors.NewServerError(err)
	}

	return user, nil
}

// RegisterUser registers given user model into the database.
// The returned user has added database auth (auth.Db) with whatever it had before.
// Returns UsernameInUseError if username is already in use.
func RegisterUser(user *models.User) (*models.User, fverrors.Error) {
	if err := models.ValidateUsername(user.Username); err != nil {
		return nil, err
	}

	_, err := dbConn.Exec(context.Background(), `INSERT INTO users (username, email, user_id) VALUES ($1, $2, $3)`,
		user.Username, user.Email, user.UserID)

	if err == nil {
		user.AuthState |= auth.Db
		return user, nil
	}

	if pgErr, okay := errors.AsType[*pgconn.PgError](err); okay && (pgErr.Code == pgerrcode.UniqueViolation) {
		return nil, fverrors.UsernameInUseError
	}

	return nil, fverrors.NewServerError(err)
}

// GetUserByUsername returns a models.User with given username field and its corresponding userID field set and
// with no auth (auth.NoAccess). Returns UserNotFoundError if no such user exists.
func GetUserByUsername(username string) (*models.User, fverrors.Error) {
	if err := models.ValidateUsername(username); err != nil {
		return nil, err
	}

	user := &models.User{
		Username:  username,
		AuthState: auth.NoAccess,
	}

	err := dbConn.QueryRow(context.Background(), `SELECT user_id FROM users WHERE username=$1`,
		username).Scan(&user.UserID)
	if (err != nil) && (errors.Is(err, pgx.ErrNoRows)) {
		return nil, fverrors.UserNotFoundError
	} else if err != nil {
		return nil, fverrors.NewServerError(err)
	}

	return user, nil
}
