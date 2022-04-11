package test

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDatabase() *pgxpool.Pool {

	user := os.Getenv("TEST_DB_USER")
	if user == "" {
		user = "postgres"
	}

	name := os.Getenv("TEST_DB_NAME")
	if name == "" {
		name = "postgres_test"
	}

	password := os.Getenv("TEST_DB_NAME")
	if password == "" {
		password = "password"
	}

	host := os.Getenv("TEST_DB_HOST")
	if host == "" {
		host = "localhost"
	}

	conn, err := pgxpool.Connect(context.Background(),
		fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable",
			user, name, password, host))
	if err != nil {
		panic("Unable to connect to database: " + err.Error())
	}

	return conn
}

func CloseDb(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), `
		DROP TABLE if exists room cascade;
	`)

	if err != nil {
		panic("Error dropping room table: " + err.Error())
	}

	_, err = db.Exec(context.Background(), `
		DROP TABLE if exists reservation cascade;
	`)

	if err != nil {
		panic("Error dropping reservation table: " + err.Error())
	}

	db.Close()
}
