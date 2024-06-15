package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/securecookie"
)

const COOKIE_NAME = "jwt"

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(secretKey []byte, userID string, expDuration time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GenerateCookie(cookieSecret *securecookie.SecureCookie, name string, value string, expDuration time.Duration, secure bool) *http.Cookie {
    encoded, err := cookieSecret.Encode(name, value)
	if err == nil {
		return &http.Cookie {
			Name:  name,
			Value: encoded,
			Path:  "/",
            MaxAge: int(expDuration.Seconds()),
            Expires: time.Now().Add(expDuration),
            Secure: secure,
            HttpOnly: secure,
		}
	}

    return nil
}

func GenerateExpiredCookie(name string) *http.Cookie {
    return &http.Cookie {
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
}

func ParseCookie(r *http.Request, cookieSecret *securecookie.SecureCookie, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	var value string
	err = cookieSecret.Decode(name, cookie.Value, &value)
	if err != nil {
		return "", err
	}

	return value, nil
}

