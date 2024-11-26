package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Open(
	user string,
	name string,
	password string,
	host string,
	port int,
) (*sqlx.DB, error) {

	dsn := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%d sslmode=disable",
		user,
		name,
		password,
		host,
		port,
	)

	return sqlx.Open("postgres", dsn)
}
