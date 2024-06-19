package models

import "github.com/Anand-S23/complete-auth/pkg/auth"

type RegisterDto struct {
    Email     string
    FirstName string
    LastName  string
    Password  string
    Confirm   string
}

type LoginDto struct {
    Email    string
    Password string
}

type OAuthRegisterDto struct {
    Provider  auth.Provider
    OAuthID   string
    Email     string
    FirstName string
    LastName  string
}

