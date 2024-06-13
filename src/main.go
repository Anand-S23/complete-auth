package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Anand-S23/complete-auth/src/auth"
	"github.com/Anand-S23/complete-auth/src/controller"
	"github.com/Anand-S23/complete-auth/src/database"
	"github.com/Anand-S23/complete-auth/src/router"
	"github.com/Anand-S23/complete-auth/src/store"
)

const (
    JwtSecretKey string = "JWT_SECRET_KEY"
    CookieHashKey string = "COOKIE_HASH_KEY"
    CookieBlockKey string = "COOKIE_BLOCK_KEY"
    DbUri string = "postgres://postgres:Password1234@postgres/complete_auth_db?sslmode=disable"
    Production bool = false
    Callback string = "http://localhost:8080/auth/google/callback"
)

func GetEnvOrPanic(env string) string {
	variable := os.Getenv(env)
	if variable == "" {
        message := fmt.Sprintf("Must provide %s variable in .env file", env)
        log.Fatal(message)
	}

	return variable
}

func main() {
    googleClientId := GetEnvOrPanic("GOOGLE_CLIENT_ID")
    googleClientSecret := GetEnvOrPanic("GOOGLE_CLIENT_SECRET")

    db := database.InitDB(DbUri, Production)
    store := store.NewStore(store.NewPgUserRepo(db))

    controller := controller.NewController(
        store, 
        []byte(JwtSecretKey), 
        []byte(CookieHashKey), 
        []byte(CookieBlockKey), 
        Production, 
        auth.NewGoogleOAuthConfig(Callback, googleClientId, googleClientSecret),
    )

    log.Println("Server running on port :8080")
    http.ListenAndServe(":8080", router.NewRouter(controller))
}

