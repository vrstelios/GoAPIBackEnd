package graphql

import (
	"GoAPIBackEnd/internal/models"
	"github.com/graphql-go/graphql"
	"strings"
)

// Create GraphQL Object type with name TaskType
var TaskType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Task",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"done": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"tasks": &graphql.Field{
			Type: graphql.NewList(TaskType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var tasks []models.Task
				for _, t := range models.LibTasks {
					tasks = append(tasks, *t)
				}
				return tasks, nil
			},
		},
		"task": &graphql.Field{
			Type: graphql.NewList(TaskType),
			Args: graphql.FieldConfigArgument{
				"done": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
				"titleContains": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"limit": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"offset": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				doneFilter, _ := params.Args["done"].(bool)
				titleContains, _ := params.Args["titleContains"].(string)
				limit, limitProvided := params.Args["limit"].(int)
				offset, offsetProvided := params.Args["offset"].(int)

				var filtered []models.Task
				for _, t := range models.LibTasks {
					if titleContains != "" && !strings.Contains(strings.ToLower(t.Title), strings.ToLower(titleContains)) {
						continue
					}
					if _, ok := params.Args["done"]; ok && t.Done != doneFilter {
						continue
					}
					filtered = append(filtered, *t)
				}

				if offsetProvided {
					if offset > len(filtered) {
						filtered = []models.Task{}
					} else {
						filtered = filtered[offset:]
					}
				}

				if limitProvided && limit < len(filtered) {
					filtered = filtered[:limit]
				}

				return filtered, nil
			},
		},
	},
})

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: RootQuery,
})
