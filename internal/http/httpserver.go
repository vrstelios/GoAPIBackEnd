package http

import (
	_ "GoAPIBackEnd/docs"
	"GoAPIBackEnd/internal/api"
	"GoAPIBackEnd/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine) {
	// Swagger documentation http://localhost:8080/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// GraphQL handler
	/*graphqlHandler := handler.New(&handler.Config{
		Schema:   &graphql.Schema,
		Pretty:   true,
		GraphiQL: true,
	})*/

	routerEndpoints := router.Group("/api")

	routerEndpoints.POST("/auth/signup", api.Signup)
	routerEndpoints.POST("/auth/login", api.Login)
	routerEndpoints.GET("/auth/logout", api.Logout)

	routerEndpoints.Use(auth.AuthJWTMiddleware())
	{
		// Exercises endpoints
		routerEndpoints.GET("/exercises", api.QueryExercises)
		routerEndpoints.GET("/exercises/:id", api.GetExercise)
		routerEndpoints.POST("/exercises", auth.RoleMiddleware("admin"), api.PostExercise)

		// Coach endpoints
		routerEndpoints.GET("/coach/:name", auth.RoleMiddleware("admin"), api.GetCoach)
		routerEndpoints.POST("/coach", auth.RoleMiddleware("admin"), api.PostCoach)

		// GraphQL endpoint
		/*routerEndpoints.POST("/graphql", func(ctx *gin.Context) {
			graphqlHandler.ServeHTTP(ctx.Writer, ctx.Request)
		})


		routerEndpoints.GET("/tasks", Query)
		routerEndpoints.GET("/tasks/:id", Get)
		routerEndpoints.POST("/tasks", RoleMiddleware("admin"), Post)
		routerEndpoints.PUT("/tasks/:id", Put)
		routerEndpoints.DELETE("/tasks/:id", Del)
		routerEndpoints.POST("/download/images", DownloadUrls)
		routerEndpoints.POST("/tasks/query", QueryTasksV2)*/
	}
}
