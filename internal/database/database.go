package database

import (
	"GoAPIBackEnd/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func InitDatabase() {
	cfg := config.GetConfig()

	var err error
	cnnDB := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Tehran",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.SSLMode)

	Conn, err = pgx.Connect(context.Background(), cnnDB)
	if err != nil {
		panic("Failed to connect to db: " + err.Error())
	}
}
