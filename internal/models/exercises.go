package models

import "time"

type Exercises struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Category    *string   `json:"category,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}
