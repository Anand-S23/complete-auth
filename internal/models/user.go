package models

import (
	"time"
)


type User struct {
    ID            string    `json:"id"`
    Email         string    `json:"email"`
    Password      string    `json:"-"`
    OAuthProvider string    `json:"oAuthProvider"`
    OAuthID       string    `json:"oAuthId"`
    LastLogin     time.Time `json:"lastLogin"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewUser(userData RegisterDto) User {
    // now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
    return User {
    }
}

