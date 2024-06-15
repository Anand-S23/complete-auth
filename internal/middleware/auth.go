package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Anand-S23/complete-auth/pkg/auth"
	"github.com/Anand-S23/complete-auth/internal/controller"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/securecookie"
)

func getUserFromResquest(r *http.Request, jwtSecretKey []byte, cookieSecret *securecookie.SecureCookie) (*auth.Claims, error) {
    tokenString, err := auth.ParseCookie(r, cookieSecret, auth.COOKIE_NAME)
	if err != nil {
        errMsg := fmt.Sprintf("invalid request, could not parse cookie: %s", err)
        return nil, errors.New(errMsg)
	}

	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
        errMsg := fmt.Sprintf("invalid cookie, could not parse token: %s", err)
        return nil, errors.New(errMsg)
	}
    
	claims, ok := token.Claims.(*auth.Claims)
	if !ok || !token.Valid {
        return nil, errors.New("invalid token, not able to parse claims")
	}

    return claims, nil
}

func Auth(next http.Handler, c *controller.Controller) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        claims, err := getUserFromResquest(r, c.JwtSecretKey, c.CookieSecret)
        if err != nil {
            log.Println(err.Error())
            controller.WriteJSON(w, http.StatusUnauthorized, controller.ErrMsg("Unauthorized"))
            return
        }
        
        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
    })
}

