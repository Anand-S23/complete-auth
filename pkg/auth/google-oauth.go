package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// TODO: make scopes dynamic and not hard coded
func NewGoogleOAuthConfig(callback string, clientId string, clientSecret string) *oauth2.Config{
    return &oauth2.Config {
        RedirectURL: callback,
        ClientID: clientId,
        ClientSecret: clientSecret,
        Scopes: []string{"https://www.googleapis.com/auth/userinfo.email"},
        Endpoint: google.Endpoint,
    }
}

