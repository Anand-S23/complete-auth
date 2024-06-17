package models

import (
	"time"

	"github.com/Anand-S23/complete-auth/pkg/auth"
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

    Profile       UserProfile
}

type UserProfile struct {
    UserID      string    `json:"userId"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	PhoneNumber string    `json:"phoneNumber"`
	PfpURL      string    `json:"pfpURL"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewUserProfile(userId string, firstName string, lastName string) UserProfile {
    return UserProfile{
        UserID: userId,
        FirstName: firstName,
        LastName: lastName,
    }
}

func NewUser(userData RegisterDto) User {
    now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

    return User {
        ID: "TODO",
        Email: userData.Email,
        Password: userData.Password,
        CreatedAt: now,
        UpdatedAt: now,
        Profile: NewUserProfile("TODO", userData.FirstName, userData.LastName),
    }
}

func NewOAuthUser(provider auth.Provider, oAuthId string, userData RegisterDto) User {
    now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

    return User {
        ID: "TODO",
        OAuthProvider: string(provider),
        OAuthID: oAuthId,
        CreatedAt: now,
        UpdatedAt: now,
        Profile: NewUserProfile("TODO", userData.FirstName, userData.LastName),
    }
}

