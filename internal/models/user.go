package models

import (
	"log"
	"strings"
	"time"

	"github.com/Anand-S23/complete-auth/pkg/auth"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
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

func (u *User) SetHashedPassword(password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
        return err
	}

    u.Password = string(hashedPassword)
    return nil
}

func (u *User) ValidatePassword(password string) error {
     return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func NewUserProfile(userId string, firstName string, lastName string, phoneNumber string, pfpURL string) UserProfile {
    return UserProfile{
        UserID: userId,
        FirstName: firstName,
        LastName: lastName,
        PhoneNumber: phoneNumber,
        PfpURL: pfpURL,
    }
}

func NewUser(userData RegisterDto) User {
    id := ulid.Make().String()

    return User {
        ID: id,
        Email: userData.Email,
        Password: userData.Password,
        Profile: NewUserProfile(id, userData.FirstName, userData.LastName, "", ""),
    }
}

func splitName(name string) (string, string) {
    arrName := strings.Split(name, " ")
    if len(arrName) == 1 {
        return name, ""
    }

    return arrName[0], arrName[len(arrName) - 1]
}

func NewOAuthUser(provider auth.Provider, userData OAuthRegisterDto) User {
    id := ulid.Make().String()
    firstName, lastName := splitName(userData.Name)

    return User {
        ID: id,
        OAuthProvider: string(provider),
        OAuthID: userData.OAuthID,
        Profile: NewUserProfile(id, firstName, lastName, "", userData.PfpURL),
    }
}

