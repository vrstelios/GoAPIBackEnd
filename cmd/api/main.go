package main

import (
	"GoAPIBackEnd/config"
	"GoAPIBackEnd/internal/model"
	"GoAPIBackEnd/internal/server"
	"fmt"
	"os"
)

// @title 		  EndPoints
// @version		  1.2
// @description   A Tag service API in Go using Gin framework
// @contact.name  DoctorVeRossi
// @contact.url   https://github.com/vrstelios/GoAPIBackEnd
// @BasePath      /api
func main() {
	model.LibTasks = make(map[string]*model.Task)

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
	//database.ConnectToDB()
	//database.SyncDatabase()
}
