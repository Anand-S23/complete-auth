package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Anand-S23/complete-auth/internal/models"
	"github.com/Anand-S23/complete-auth/pkg/auth"
)

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) error {
    var userData models.RegisterDto
    err := json.NewDecoder(r.Body).Decode(&userData)
    if err != nil {
        log.Println("error parsing register data")
        return WriteJSON(w, http.StatusBadRequest, ErrMsg("Could not parse register data"))
    }

    // TODO: validate user data - make sure is valid and user does not exist

    user := models.NewUser(userData)
    err = user.SetHashedPassword(userData.Password)
    if err != nil {
        log.Println("error hashing password")
        return WriteJSON(w, http.StatusInternalServerError, ErrMsg("Internal server error occured, please try again later"))
    }

    err = c.store.UserRepo.InsertUser(context.Background(), &user)
	if err != nil {
        log.Println("error inserting user into db")
        return WriteJSON(w, http.StatusInternalServerError, ErrMsg("Internal server error occured, please try again later"))
	}

    successMsg := map[string]string {
        "message": "User created successfully",
        "userID": user.ID,
    }
    log.Println("User created successfully")
    return WriteJSON(w, http.StatusOK, successMsg)
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) error {
    var loginData models.LoginDto
    err := json.NewDecoder(r.Body).Decode(&loginData)
    if err != nil {
        log.Println("error parsing login data")
        return WriteJSON(w, http.StatusBadRequest, ErrMsg("Could not parse login data"))
    }

    user, err := c.store.UserRepo.GetUserByEmail(context.Background(), loginData.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Println("user trying to sign in with non existing email")
        } else {
            log.Println("user entered wrong password, could not sign in")
        }
        return WriteJSON(w, http.StatusBadRequest, ErrMsg("Incorrect email or password, please try again"))
    }

    err = user.ValidatePassword(loginData.Password)
	if err != nil {
        log.Println("Passwords do not match")
        return WriteJSON(w, http.StatusBadRequest, ErrMsg("Incorrect email or password, please try again"))
	}

    expDuration := time.Hour * 24
    token, err := auth.GenerateToken(c.JwtSecretKey, user.ID, expDuration)
    if err != nil {
        log.Println("Error generating token")
        return WriteJSON(w, http.StatusInternalServerError, ErrMsg("Internal server error occured, please try again later"))
    }

    cookie := auth.GenerateCookie(c.CookieSecret, auth.COOKIE_NAME, token, expDuration, c.production)
    if cookie == nil {
        log.Println("Error generating cookie")
        return WriteJSON(w, http.StatusInternalServerError, ErrMsg("Internal server error occured, please try again later"))
    }
    http.SetCookie(w, cookie)

    successMsg := map[string]string {
        "message": "User logged in successfully",
    }
    log.Println("User successfully logged in")
    c.store.UserRepo.UpdateLastLogin(context.Background(), user.ID)
    return WriteJSON(w, http.StatusOK, successMsg)
}

func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) error {
    cookie := auth.GenerateExpiredCookie(auth.COOKIE_NAME)
    http.SetCookie(w, cookie)
    log.Println("User successfully logged out")
    return WriteJSON(w, http.StatusOK, "")
}

func (c *Controller) GetAuthUserID(w http.ResponseWriter, r *http.Request) error {
    currentUserID := r.Context().Value("user_id").(string)
    return WriteJSON(w, http.StatusOK, currentUserID)
}

