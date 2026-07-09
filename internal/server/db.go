package server

import (
	"context"
	"fmt"
	"os"

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
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (s *Server) registerUser(user *User) error {
	_, err := s.PgConn.Exec(context.Background(),
		"INSERT INTO users (userName, displayName, email, userID) values ($1, $2, $3, $4)",
		user.Username, user.Name, user.Email, user.GoogleUserID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) getUser(u *User) error {
	ok := validateUserID(u.GoogleUserID)
	if !ok {
		return fmt.Errorf("invalid user id")
	}

	err := s.PgConn.QueryRow(context.Background(), "SELECT userName, displayName FROM users WHERE id=$1", u.GoogleUserID).Scan(&u.Username, &u.Name)
	if (err != nil) && (err == pgx.ErrNoRows) {
		return err
	}

	return nil
}
