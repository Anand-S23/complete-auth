package controller

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Anand-S23/complete-auth/internal/models"
	"github.com/Anand-S23/complete-auth/pkg/auth"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
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
    err = c.store.UserRepo.InsertUser(context.Background(), &user)
	if err != nil {
        log.Println("Error inserting user into db")
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
        return WriteJSON(w, http.StatusBadRequest, ErrMsg("Could not parse login data"))
    }

    user, err := c.store.UserRepo.GetUserByEmail(context.Background(), loginData.Email)
    if err != nil {
        // TODO: figure out what error this is because it could be an actual password and not the user typing in wrong email
        log.Println("Passwords do not match")
        return WriteJSON(w, http.StatusBadRequest, ErrMsg("Incorrect email or password, please try again"))
    }

    // TODO: Maybe move abstract this out
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

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func generateStateOAuthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(code string, config *oauth2.Config) ([]byte, error) {
	// Use code to get token and get user info from Google.

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

func (c *Controller) GoogleLogin(w http.ResponseWriter, r *http.Request) error {
	oauthState := generateStateOAuthCookie(w)
	u := c.googleOAuthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
    return nil
}

func (c *Controller) GoogleCallback(w http.ResponseWriter, r *http.Request) error {
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil
	}

	data, err := getUserDataFromGoogle(r.FormValue("code"), c.googleOAuthConfig)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil
	}

    var oauthRegisterData models.OAuthRegisterDto
    err = json.Unmarshal(data, &oauthRegisterData)
    if err != nil {
        // TODO: this might not work as expected, might need to acutally redirect here
        return WriteJSON(w, http.StatusBadRequest, ErrMsg("Could not parse oauth data"))
    }

    user, err := c.store.UserRepo.GetBaseUserByEmail(context.Background(), oauthRegisterData.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            oauthUser := models.NewOAuthUser(auth.ProviderGoogle, oauthRegisterData)
            err = c.store.UserRepo.InsertUser(context.Background(), &oauthUser)
            if err != nil {
                // TODO: this might not work as expected, might need to acutally redirect here
                return WriteJSON(w, http.StatusBadRequest, ErrMsg("Could not insert user into database"))
            }
        } else {
            // TODO: this might not work as expected, might need to acutally redirect here
            return WriteJSON(w, http.StatusBadRequest, ErrMsg("Could not get user from database"))
        }
    }
    c.store.UserRepo.UpdateLastLogin(context.Background(), user.ID)
	// Redirect or response with a token.

	fmt.Fprintf(w, "UserInfo: %s\n", data)
    return nil
}

