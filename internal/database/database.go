package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"os"
)

var Conn *pgx.Conn

func InitDatabase() {
	var err error
	Conn, err = pgx.Connect(context.Background(), os.Getenv("DB"))
	if err != nil {
		panic("Failed to connect to db: " + err.Error())
	}
}
