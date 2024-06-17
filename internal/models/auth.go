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

