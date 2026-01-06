package api

import (
	"GoAPIBackEnd/config"
	"GoAPIBackEnd/internal/type/models"
)

var (
	DefaultClient *APIClient
)

func Init() {
	DefaultClient = NewAPIClient(config.GetConfig().API.BaseURL)
}

func ClientQueryWorkouts(userId string) ([]models.Workouts, error) {
	if DefaultClient == nil {
		Init()
	}
	return DefaultClient.QueryWorkouts(userId)
}

func ClientCreateWorkoutLogs(logs []models.WorkoutLog) ([]models.WorkoutLog, error) {
	if DefaultClient == nil {
		Init()
	}
	return DefaultClient.CreateWorkoutLogs(logs)
}
