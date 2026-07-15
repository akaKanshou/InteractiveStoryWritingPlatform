package server

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
)

func Connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v",
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_DB"),
	))

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to DB")
	return conn, nil
}

func validateUserID(id string) bool {
	for _, c := range id {
		if (c < '0') || (c > '9') {
			return false
		}
	}
	return true
}

func validateUsername(name string) error {
	name = strings.ToLower(name)

	_, err := regexp.Match("^[a-z0-9_]{3,18}", []byte(name))

	return err
}

func (s *Server) registerUser(user *user) error {
	_, err := s.PgConn.Exec(context.Background(),
		`INSERT INTO users (username, email, user_id) values ($1, $2, $3)`,
		user.Username, user.Email, user.GoogleUserID)

	if err != nil {
		return err
	}

	user.AuthState |= AuthDB
	return nil
}

func (s *Server) getUser(u *user) error {
	ok := validateUserID(u.GoogleUserID)

	if !ok {
		fmt.Println("Invalid user ID")
		return fmt.Errorf("invalid user id")
	}

	err := s.PgConn.QueryRow(context.Background(), `SELECT username FROM users WHERE user_id=$1`, u.GoogleUserID).Scan(&u.Username)
	if (err != nil) && (err == pgx.ErrNoRows) {
		return err
	} else if err != nil {
		fmt.Println("Error getting user:", u.GoogleUserID, err)
		return err
	}

	u.AuthState |= AuthDB
	return nil
}
