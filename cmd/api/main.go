package main

import (
	"GoAPIBackEnd/internal"
	"GoAPIBackEnd/internal/api"
	"GoAPIBackEnd/internal/database"
	"fmt"
	"os"
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

	srv.Run(os.Getenv("PORT"))
}

func init() {
	database.InitDatabase()
	api.Init()
}
