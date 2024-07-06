package models

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
    OAuthID       string `json:"id"`
    Email         string `json:"email"`
    VerifiedEmail bool   `json:"verified_email"`
    Name          string `json:"name"`
    PfpURL        string `json:"picture"`
}

