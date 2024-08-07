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
	"golang.org/x/oauth2"
)

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
    baseRedirectURL := fmt.Sprintf("%s/", c.feURI)
    serverErrorRedirectURL := fmt.Sprintf("%s/login?serverError", c.feURI)
    accountTypeErrorRedirectURL := fmt.Sprintf("%s/login?accountTypeError", c.feURI)

	oauthState, _ := r.Cookie("oauthstate")
	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, serverErrorRedirectURL, http.StatusTemporaryRedirect)
		return nil
	}

	data, err := getUserDataFromGoogle(r.FormValue("code"), c.googleOAuthConfig)
	if err != nil {
		log.Println("could not get user data from google")
		http.Redirect(w, r, serverErrorRedirectURL, http.StatusTemporaryRedirect)
		return nil
	}

    var oauthRegisterData models.OAuthRegisterDto
    err = json.Unmarshal(data, &oauthRegisterData)
    if err != nil {
		log.Println("error serializing register data")
		http.Redirect(w, r, serverErrorRedirectURL, http.StatusTemporaryRedirect)
		return nil
    }

    user, err := c.store.UserRepo.GetBaseUserByEmail(context.Background(), oauthRegisterData.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Println("user does not exist creating user")
            oauthUser := models.NewOAuthUser(auth.ProviderGoogle, oauthRegisterData)
            err = c.store.UserRepo.InsertUser(context.Background(), &oauthUser)
            if err != nil {
                log.Println("error inserting user into db")
                http.Redirect(w, r, serverErrorRedirectURL, http.StatusTemporaryRedirect)
                return nil
            }
        } else {
            log.Println("error reteriving user into db")
            http.Redirect(w, r, serverErrorRedirectURL, http.StatusTemporaryRedirect)
            return nil
        }
    } else if user.OAuthProvider != string(auth.ProviderGoogle) {
        log.Println("error reteriving user into db")
        http.Redirect(w, r, accountTypeErrorRedirectURL, http.StatusTemporaryRedirect)
        return nil
    }

    expDuration := time.Hour * 24
    token, err := auth.GenerateToken(c.JwtSecretKey, user.ID, expDuration)
    if err != nil {
        log.Println("Error generating token")
    }

    cookie := auth.GenerateCookie(c.CookieSecret, auth.COOKIE_NAME, token, expDuration, c.production)
    if cookie == nil {
        log.Println("Error generating cookie")
    }
    http.SetCookie(w, cookie)

	// fmt.Fprintf(w, "UserInfo: %s\n", data)
    log.Println("login sucessful")
    c.store.UserRepo.UpdateLastLogin(context.Background(), user.ID)
    http.Redirect(w, r, baseRedirectURL, http.StatusTemporaryRedirect)
    return nil
}

