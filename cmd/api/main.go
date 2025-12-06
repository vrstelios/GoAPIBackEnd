package main

import (
	"GoAPIBackEnd/config"
	"GoAPIBackEnd/internal/database"
	"GoAPIBackEnd/internal/models"
	"GoAPIBackEnd/internal/server"
	"fmt"
	"os"
)

// @title 		  EndPoints
// @version		  1.3
// @description   A Tag service API in Go using Gin framework
// @contact.name  DoctorVeRossi
// @contact.url   https://github.com/vrstelios/GoAPIBackEnd
// @BasePath      /api
func main() {
	models.LibTasks = make(map[string]*models.Task)

	srv := server.New()

	fmt.Println(`
	 ______     ______         ______     ______   __
	/\  ___\   /\  __ \       /\  __ \   /\  == \ /\ \
	\ \ \__ \  \ \ \/\ \   -  \ \  __ \  \ \  _-/ \ \ \
	 \ \_____\  \ \_____\  -   \ \_\ \_\  \ \_\    \ \_\
	  \/_____/   \/_____/       \/_/\/_/   \/_/     \/_/ `)

	srv.Run(os.Getenv("PORT"))
}

func init() {
	config.LoadEnvVariables()
	database.ConnectToDB()
	database.SyncDatabase()
}
