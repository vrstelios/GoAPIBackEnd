package models

import "time"

type Workouts struct {
	Id          string             `json:"id"`
	UserId      string             `json:"userId"`
	Name        string             `json:"name"`
	Notes       *string            `json:"notes,omitempty"`
	ScheduledAt *time.Time         `json:"scheduledAt,omitempty"`
	CreatedAt   time.Time          `json:"createdAt"`
	Relations   *WorkoutsRelations `json:"relationships,omitempty"`
}

type WorkoutsRelations struct {
	Users *Users `json:"Users"`
}

type WorkoutExercises struct {
	Id            string              `json:"id"`
	WorkoutId     string              `json:"workoutId"`
	ExerciseId    string              `json:"exerciseId"`
	Sets          *int                `json:"sets,omitempty"`
	Reps          *int                `json:"reps,omitempty"`
	Weight        *float64            `json:"weight,omitempty"`
	RestSeconds   *string             `json:"restSeconds,omitempty"`
	ExerciseOrder int                 `json:"exerciseOrder"`
	Relations     *WorkoutExRelations `json:"relationships,omitempty"`
}

type WorkoutExRelations struct {
	Workout  *Workouts  `json:"workout"`
	Exercise *Exercises `json:"exercise"`
}

type WorkoutLog struct {
	Id          string               `json:"id"`
	WorkoutId   string               `json:"workoutId"`
	UserId      string               `json:"userId"`
	PerformedAt *time.Time           `json:"performedAt,omitempty"`
	Duration    *float64             `json:"duration,omitempty"`
	TotalVolume *string              `json:"totalVolume,omitempty"`
	Notes       *string              `json:"notes,omitempty"`
	Relations   *WorkoutLogRelations `json:"relationships,omitempty"`
}

type WorkoutLogRelations struct {
	Workout *Workouts `json:"workout"`
	User    *Users    `json:"user"`
}

type LoadWorkouts struct {
	Excel    string `json:"excel"`
	Success  bool   `json:"success"`
	Inserted int    `json:"inserted"`
	Failed   int    `json:"failed"`
	Error    error  `json:"error"`
}
