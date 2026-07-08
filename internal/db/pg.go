package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
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
