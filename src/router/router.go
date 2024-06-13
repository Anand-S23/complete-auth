package router

import (
	"net/http"

	"github.com/Anand-S23/complete-auth/src/controller"
	"github.com/Anand-S23/complete-auth/src/middleware"
	"github.com/gorilla/handlers"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

func Fn(fn apiFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        err := fn(w, r)
        if err != nil {
            controller.WriteJSON(w, http.StatusInternalServerError, controller.ErrMsg(err.Error()))
        }
    }
}

func NewRouter(c *controller.Controller) http.Handler {
    router := http.NewServeMux()
    router.HandleFunc("GET /ping", Fn(c.Ping)) // Health Check

    router.HandleFunc("POST /register", Fn(c.Register))
    router.HandleFunc("POST /login", Fn(c.Login))
    router.HandleFunc("POST /logout", Fn(c.Logout))
    router.HandleFunc("GET /getAuthUserID", middleware.Auth(Fn(c.GetAuthUserID), c))

    router.HandleFunc("GET /auth/google/login", Fn(c.GoogleLogin))
    router.HandleFunc("GET /auth/google/callback", Fn(c.GoogleCallback))

    corsHandler := handlers.CORS(
        handlers.AllowedOrigins([]string{"http://localhost:5713"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

    return corsHandler(router)
}

