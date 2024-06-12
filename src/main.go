package main

import (
	"log"
	"net/http"

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
)

func main() {
    db := database.InitDB(DbUri, Production)
    store := store.NewStore(store.NewPgUserRepo(db))
    controller := controller.NewController(
        store, []byte(JwtSecretKey), []byte(CookieHashKey), []byte(CookieBlockKey), Production)

    log.Println("Server running on port :8080")
    http.ListenAndServe(":8080", router.NewRouter(controller))
}

