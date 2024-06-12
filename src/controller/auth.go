package controller

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Anand-S23/complete-auth/src/auth"
	"github.com/Anand-S23/complete-auth/src/models"
	"golang.org/x/crypto/bcrypt"
)

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) error {
    var userData models.RegisterDto
    err := json.NewDecoder(r.Body).Decode(&userData)
    if err != nil {
        log.Println("Error parsing register data")
        return WriteJSON(w, http.StatusBadRequest, ErrMsg("Could not parse register data"))
    }

    // TODO: validate user data - make sure is valid and user does not exist

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
        log.Println("Error hashing the password")
        return WriteJSON(w, http.StatusInternalServerError, ErrMsg("Internal server error occured, please try again later"))
	}
    userData.Password = string(hashedPassword)

    user := models.NewUser(userData)
    // TODO: Save to db

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
        return WriteJSON(w, http.StatusBadRequest, ErrMsg("Could not parse login data"))
    }

    // TODO: Try to get user from db
    
    user := models.User{}
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
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

