package main

import (
	"GoAPIBackEnd/config"
	"GoAPIBackEnd/internal"
	"GoAPIBackEnd/internal/api"
	config2 "GoAPIBackEnd/internal/config"
	"fmt"
)

// @title 		  EndPoints
// @version		  1.7
// @description   A Tag service API in Go using Gin framework
// @contact.name  DoctorVeRossi
// @contact.url   https://github.com/vrstelios/GoAPIBackEnd
// @BasePath      /api
func main() {
	srv := internal.New()

	fmt.Println(`
	 ______     ______         ______     ______   __
	/\  ___\   /\  __ \       /\  __ \   /\  == \ /\ \
	\ \ \__ \  \ \ \/\ \   -  \ \  __ \  \ \  _-/ \ \ \
	 \ \_____\  \ \_____\  -   \ \_\ \_\  \ \_\    \ \_\
	  \/_____/   \/_____/       \/_/\/_/   \/_/     \/_/ `)

	srv.Run(config.GetConfig().Server.Port)
}

func init() {
	config2.InitDatabase()
	api.Init()
}
