package http

import (
	_ "GoAPIBackEnd/docs"
	"GoAPIBackEnd/internal/api"
	"GoAPIBackEnd/internal/auth"
	"GoAPIBackEnd/internal/graphql"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine) {
	// Swagger documentation http://localhost:8080/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// GraphQL handler
	graphqlHandler := handler.New(&handler.Config{
		Schema:   &graphql.Schema,
		Pretty:   true,
		GraphiQL: true,
	})

	routerEndpoints := router.Group("/api")

	routerEndpoints.POST("/auth/signup", api.Signup)
	routerEndpoints.POST("/auth/login", api.Login)
	routerEndpoints.GET("/auth/logout", api.Logout)

	// GraphQL endpoint
	// Found workouts with userId and created logs
	// Combine API Gateway, joins, filters, pagination, still and complex reports, in any query
	routerEndpoints.POST("/graphql/workouts/logs", func(ctx *gin.Context) {
		graphqlHandler.ServeHTTP(ctx.Writer, ctx.Request)
	})
	routerEndpoints.POST("/workoutLogs", api.PostWorkoutLogs)
	routerEndpoints.GET("/workouts", api.QueryWorkouts)

	routerEndpoints.Use(auth.AuthJWTMiddleware())
	{
		// Exercises endpoints
		routerEndpoints.GET("/exercises", api.QueryExercises)
		routerEndpoints.GET("/exercises/:id", api.GetExercise)
		routerEndpoints.POST("/exercises", auth.RoleMiddleware("admin"), api.PostExercise)

		// Coach endpoints
		routerEndpoints.GET("/coach/:name", auth.RoleMiddleware("admin"), api.GetCoach)
		routerEndpoints.POST("/coach", auth.RoleMiddleware("admin"), api.PostCoach)

		// Workout endpoints
		routerEndpoints.POST("/workouts", api.PostWorkout)
		routerEndpoints.GET("/workouts/:id", api.GetWorkout)
		//routerEndpoints.GET("/workouts", api.QueryWorkouts)
		routerEndpoints.PUT("/workouts/:id", api.PutWorkout)
		routerEndpoints.DELETE("/workouts/:id", api.DelWorkout)
		routerEndpoints.POST("/workouts/query", api.QueryWorkoutsV2)

		// WorkoutLog endpoints
		//routerEndpoints.POST("/workoutLogs", api.PostWorkoutLogs)

		// Routing&Channel
		//routerEndpoints.POST("/download/images", DownloadUrls) xlsx
		//routerEndpoints.POST("/tasks/query", QueryTasksV2)*/
	}
}
