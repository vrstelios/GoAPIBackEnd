package models

import (
	"github.com/google/uuid"
)

type Users struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
}

type Login struct {
	HashedPassword string
	SessionToken   string
	Role           string
}
