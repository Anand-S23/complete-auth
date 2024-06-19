package models

import (
	"time"

	"github.com/Anand-S23/complete-auth/pkg/auth"
	"github.com/oklog/ulid/v2"
)

type User struct {
    ID            string    `db:"id"             json:"id"`
    Email         string    `db:"email"          json:"email"`
    Password      string    `db:"password_hash"  json:"-"`
    OAuthProvider string    `db:"oauth_provider" json:"oAuthProvider"`
    OAuthID       string    `db:"oauth_id"       json:"oAuthId"`
    LastLogin     time.Time `db:"last_login"     json:"lastLogin"`
	CreatedAt     time.Time `db:"created_at"     json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at"     json:"updatedAt"`

    Profile       UserProfile
}

type UserProfile struct {
    UserID      string    `db:"user_id"      json:"userId"`
	FirstName   string    `db:"first_name"   json:"firstName"`
	LastName    string    `db:"last_name"    json:"lastName"`
	PhoneNumber string    `db:"phone_number" json:"phoneNumber"`
	PfpURL      string    `db:"pfp_url"      json:"pfpURL"`
	CreatedAt   time.Time `db:"created_at"   json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at"   json:"updatedAt"`
}

func NewUserProfile(userId string, firstName string, lastName string) UserProfile {
    now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

    return UserProfile{
        UserID: userId,
        FirstName: firstName,
        LastName: lastName,
        CreatedAt: now,
        UpdatedAt: now,
    }
}

func NewUser(userData RegisterDto) User {
    now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
    id := ulid.Make().String()

    return User {
        ID: id,
        Email: userData.Email,
        Password: userData.Password,
        CreatedAt: now,
        UpdatedAt: now,
        Profile: NewUserProfile(id, userData.FirstName, userData.LastName),
    }
}

func NewOAuthUser(provider auth.Provider, oAuthId string, userData RegisterDto) User {
    now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
    id := ulid.Make().String()

    return User {
        ID: id,
        OAuthProvider: string(provider),
        OAuthID: oAuthId,
        CreatedAt: now,
        UpdatedAt: now,
        Profile: NewUserProfile(id, userData.FirstName, userData.LastName),
    }
}

