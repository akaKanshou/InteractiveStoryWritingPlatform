package server

import (
	"context"
	"fmt"
	"os"
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

func validateName(field, name string, minLen, maxLen int, allowSpace bool) error {
	name = strings.ToLower(name)

	if (len(name) < minLen) || (len(name) > maxLen) {
		return fmt.Errorf("%v must be between %d and %d characters", field, minLen, maxLen)
	}

	for _, char := range name {
		if (char == ' ') && allowSpace {
			continue
		}

		if (char >= '0') && (char <= '9') {
			continue
		}

		if (char >= 'a') && (char <= 'z') {
			continue
		}

		return fmt.Errorf("%v contains invalid characters", field)
	}

	return nil
}

func (s *Server) registerUser(user *User) error {
	_, err := s.PgConn.Exec(context.Background(),
		`INSERT INTO users (username, display_name, email, user_id) values ($1, $2, $3, $4)`,
		user.Username, user.Name, user.Email, user.GoogleUserID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) getUser(u *User) error {
	ok := validateUserID(u.GoogleUserID)

	if !ok {
		fmt.Println("Invalid user ID")
		return fmt.Errorf("invalid user id")
	}

	err := s.PgConn.QueryRow(context.Background(), `SELECT username, display_name FROM users WHERE user_id=$1`, u.GoogleUserID).Scan(&u.Username, &u.Name)
	if (err != nil) && (err == pgx.ErrNoRows) {
		return err
	} else if err != nil {
		fmt.Println("Error getting user:", u.GoogleUserID, err)
		return err
	}

	return nil
}
