package models

import (
	"time"
)

type Users struct {
	Id           string          `json:"id"`
	Name         string          `json:"name"`
	Password     string          `json:"password"`
	Email        string          `json:"email"`
	Role         string          `json:"role"`
	Token        *string         `json:"token,omitempty"`
	RefreshToken *string         `json:"refreshToken,omitempty"`
	CoachId      *string         `json:"coachId"`
	CreatedAt    time.Time       `json:"createdAt"`
	Relations    *UsersRelations `json:"relationships,omitempty"`
}

type UsersRelations struct {
	Coach *Coach `json:"coach"`
}

type Coach struct {
	Id        string          `json:"id"`
	Name      string          `json:"name"`
	UserId    string          `json:"userId"`
	Relations *CoachRelations `json:"relationships,omitempty"`
}

type CoachRelations struct {
	User *Users `json:"users"`
}

type Login struct {
	HashedPassword string
	SessionToken   string
	Role           string
}
