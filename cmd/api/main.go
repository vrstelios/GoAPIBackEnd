package main

import (
	"GoAPIBackEnd/config"
	"GoAPIBackEnd/internal/model"
	"GoAPIBackEnd/internal/server"
	"fmt"
)

// @title 		  EndPoints
// @version		  1.2
// @description   A Tag service API in Go using Gin framework
// @contact.name  DoctorVeRossi
// @contact.url   https://github.com/vrstelios/GoAPIBackEnd
// @BasePath      /api
func main() {
	config.Load()
	model.LibTasks = make(map[string]*model.Task)

	srv := server.New()

	fmt.Println(`
	 ______     ______         ______     ______   __
	/\  ___\   /\  __ \       /\  __ \   /\  == \ /\ \
	\ \ \__ \  \ \ \/\ \   -  \ \  __ \  \ \  _-/ \ \ \
	 \ \_____\  \ \_____\  -   \ \_\ \_\  \ \_\    \ \_\
	  \/_____/   \/_____/       \/_/\/_/   \/_/     \/_/ `)

	srv.Run(config.App.MasterAPIDomain)
}
