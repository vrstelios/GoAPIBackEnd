package main

import (
	"GoAPIBackEnd/core"
	_ "GoAPIBackEnd/docs"
	"GoAPIBackEnd/model"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func AuthSessionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := core.Authorize(ctx); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		ctx.Next()
	}
}

// @title 		  EndPoints
// @version		  1.0
// @description   A Tag service API in Go using Gin framework
// @contact.name  DoctorVeRossi
// @contact.url   https://github.com/vrstelios/....
// @BasePath      /api
func main() {
	model.LibTasks = make(map[string]*model.Task)

	router := gin.Default()

	//http://localhost:8080/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // add swagger

	routerEndpoints := router.Group("/api")

	routerEndpoints.POST("/auth/register", core.Register)
	routerEndpoints.POST("/auth/login", core.Login)
	routerEndpoints.POST("/auth/protected", core.Protected)
	routerEndpoints.GET("/auth/logout", core.Logout)

	routerEndpoints.Use(AuthSessionMiddleware())
	{
		routerEndpoints.GET("/tasks", core.Query)
		routerEndpoints.GET("/tasks/:id", core.Get)
		routerEndpoints.POST("/tasks", core.Post)
		routerEndpoints.PUT("/tasks/:id", core.Put)
		routerEndpoints.DELETE("/tasks/:id", core.Del)
		routerEndpoints.POST("/download/images", core.DownloadUrls)
		routerEndpoints.POST("/tasks/query", core.QueryTasksV2)
	}

	fmt.Println(`
 ______     ______         ______     ______   __    
/\  ___\   /\  __ \       /\  __ \   /\  == \ /\ \   
\ \ \__ \  \ \ \/\ \   -  \ \  __ \  \ \  _-/ \ \ \  
 \ \_____\  \ \_____\  -   \ \_\ \_\  \ \_\    \ \_\ 
  \/_____/   \/_____/       \/_/\/_/   \/_/     \/_/ `)
	router.Run("localhost:8080")
}
