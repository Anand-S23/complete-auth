package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/amazon"
	"golang.org/x/oauth2/bitbucket"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/gitlab"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/linkedin"
	"golang.org/x/oauth2/microsoft"
	"golang.org/x/oauth2/paypal"
	"golang.org/x/oauth2/slack"
	"golang.org/x/oauth2/spotify"
	"golang.org/x/oauth2/twitch"
)

type Provider string

const (
    ProviderAmazon    Provider = "amazon"
    ProviderBitbucket Provider = "bitbucket"
    ProviderFacebook  Provider = "facebook"
    ProviderGithub    Provider = "github"
    ProviderGitlab    Provider = "gitlab"
    ProviderGoogle    Provider = "google"
    ProviderLinkedIn  Provider = "linkedin"
    ProviderMicrosoft Provider = "microsoft"
    ProviderPayPal    Provider = "paypal"
    ProviderSlack     Provider = "slack"
    ProviderSpotify   Provider = "spotify"
    ProviderTwitch    Provider = "twitch"
)

var ProviderEndpoints map[Provider]oauth2.Endpoint = map[Provider]oauth2.Endpoint{
    "amazon": amazon.Endpoint,
    "bitbucket": bitbucket.Endpoint,
    "facebook": facebook.Endpoint,
    "github": github.Endpoint,
    "gitlab": gitlab.Endpoint,
    "google": google.Endpoint,
    "linkedin": linkedin.Endpoint,
    "microsoft": microsoft.LiveConnectEndpoint,
    "paypal": paypal.Endpoint,
    "slack": slack.Endpoint,
    "spotify": spotify.Endpoint,
    "twitch": twitch.Endpoint,
}

func NewOAuthConfig(provider Provider, callback string, clientId string, clientSecret string, scopes []string) *oauth2.Config {
    return &oauth2.Config {
        RedirectURL: callback,
        ClientID: clientId,
        ClientSecret: clientSecret,
        Scopes: scopes,
        Endpoint: ProviderEndpoints[provider],
    }
}

