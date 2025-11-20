package server

import (
	_ "GoAPIBackEnd/docs"
	"GoAPIBackEnd/internal/graphql"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func setupRoutes(router *gin.Engine) {
	// Swagger documentation
	// http://localhost:8080/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// GraphQL handler
	graphqlHandler := handler.New(&handler.Config{
		Schema:   &graphql.Schema,
		Pretty:   true,
		GraphiQL: true,
	})

	routerEndpoints := router.Group("/api")

	routerEndpoints.POST("/auth/register", Register)
	routerEndpoints.POST("/auth/login", Login)
	routerEndpoints.GET("/auth/logout", Logout)

	routerEndpoints.Use(AuthJWTMiddleware())
	{
		// GraphQL endpoint
		routerEndpoints.POST("/graphql", func(ctx *gin.Context) {
			graphqlHandler.ServeHTTP(ctx.Writer, ctx.Request)
		})

		routerEndpoints.GET("/tasks", Query)
		routerEndpoints.GET("/tasks/:id", Get)
		routerEndpoints.POST("/tasks", RoleMiddleware("admin"), Post)
		routerEndpoints.PUT("/tasks/:id", Put)
		routerEndpoints.DELETE("/tasks/:id", Del)
		routerEndpoints.POST("/download/images", DownloadUrls)
		routerEndpoints.POST("/tasks/query", QueryTasksV2)
	}
}
