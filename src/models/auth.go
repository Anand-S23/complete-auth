package models

type RegisterDto struct {
    FirstName string
    LastName  string
    Email     string
    Phone     string
    Password  string
    Confirm   string
}

type LoginDto struct {
    Email    string
    Password string
}

