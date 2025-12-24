package graphql

import (
	"GoAPIBackEnd/internal/database"
	"GoAPIBackEnd/internal/type/models"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

// Create GraphQL Object type with name WorkoutLogType
var WorkoutLogType = graphql.NewObject(graphql.ObjectConfig{
	Name: "WorkoutLog",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"workoutId": &graphql.Field{
			Type: graphql.String,
		},
		"userId": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"completeWorkout": &graphql.Field{
			Type: WorkoutLogType,
			Args: graphql.FieldConfigArgument{
				"workoutId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"userId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var dbConn = database.Conn
				workoutId := params.Args["workoutId"].(string)
				userId := params.Args["userId"].(string)

				log := models.WorkoutLog{
					Id:        uuid.NewString(),
					WorkoutId: workoutId,
					UserId:    userId,
				}

				err := database.PostWorkoutLog(dbConn, log)
				if err != nil {
					return nil, err
				}

				return log, nil
			},
		},
	},
})

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name:   "RootQuery",
	Fields: graphql.Fields{},
})

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    RootQuery, //continue
	Mutation: RootMutation,
})
