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

