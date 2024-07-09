package main

import (
	"log"
	"net/http"

	"github.com/Anand-S23/complete-auth/internal/controller"
	"github.com/Anand-S23/complete-auth/internal/database"
	"github.com/Anand-S23/complete-auth/internal/router"
	"github.com/Anand-S23/complete-auth/internal/store"
	"github.com/Anand-S23/complete-auth/pkg/auth"
	"github.com/Anand-S23/complete-auth/pkg/config"
)

func main() {
    env := config.LoadEnv()

    db := database.InitDB(env.DB_URI, env.PRODUCTION)
    store := store.NewStore(store.NewPgUserRepo(db))

    googleOAuthConfig := auth.NewOAuthConfig(
        auth.ProviderGoogle, 
        env.GOOGLE_CALLBACK_URI, 
        env.GOOGLE_CLIENT_ID, 
        env.GOOGLE_CLIENT_SECRET,
        []string{"https://www.googleapis.com/auth/userinfo.email"},
    )
    controller := controller.NewController(store, env, googleOAuthConfig)

    log.Println("Server running on port: ", env.PORT);
    baseRouter := router.NewRouter(controller)
    router := router.NewCorsRouter(baseRouter, env.FE_URI)
    http.ListenAndServe(":" + env.PORT, router)
}

