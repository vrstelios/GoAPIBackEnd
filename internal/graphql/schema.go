package graphql

import (
	"GoAPIBackEnd/internal/api"
	"GoAPIBackEnd/internal/type/models"
	"errors"
	"fmt"
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
		"userId": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"completeWorkout": &graphql.Field{
			Type: graphql.NewList(WorkoutLogType),
			Args: graphql.FieldConfigArgument{
				"userId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userId := params.Args["userId"].(string)
				if userId == "" {
					return nil, errors.New("userId is required")
				}

				wks, err := api.ClientQueryWorkouts(userId)
				if err != nil {
					return nil, fmt.Errorf("failed to get workouts: %w", err)
				}

				logs := make([]models.WorkoutLog, 0)
				for _, wk := range wks {
					logs = append(logs, models.WorkoutLog{
						Id:        uuid.NewString(),
						WorkoutId: wk.Id,
						UserId:    userId,
					})
				}

				createdLogs, err := api.ClientCreateWorkoutLogs(logs)
				if err != nil {
					return nil, fmt.Errorf("failed to save workout logs: %w", err)
				}

				return createdLogs, nil
			},
		},
	},
})

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
})

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    RootQuery,
	Mutation: RootMutation,
})
