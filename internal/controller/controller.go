package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Anand-S23/complete-auth/internal/store"
	"github.com/Anand-S23/complete-auth/pkg/config"
	"github.com/gorilla/securecookie"
	"golang.org/x/oauth2"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func ErrMsg(message string) map[string]string {
    return map[string]string {"error": message}
}

type Controller struct {
    store *store.Store
    production   bool
    googleOAuthConfig *oauth2.Config
    feURI        string
    JwtSecretKey []byte
    CookieSecret *securecookie.SecureCookie
}

func NewController(store *store.Store, env *config.EnvVars, oauthConfig *oauth2.Config) *Controller {
    return &Controller {
        store: store,
        production: env.PRODUCTION,
        googleOAuthConfig: oauthConfig,
        feURI: env.FE_URI,
        JwtSecretKey: env.JWT_SECRET,
        CookieSecret: securecookie.New(env.COOKIE_HASH_KEY, env.COOKIE_BLOCK_KEY),
    }
}

func (c *Controller) Ping(w http.ResponseWriter, r *http.Request) error {
    return WriteJSON(w, http.StatusOK, "Pong")
}

