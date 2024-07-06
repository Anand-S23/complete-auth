package router

import (
	"net/http"

	"github.com/Anand-S23/complete-auth/internal/controller"
	"github.com/Anand-S23/complete-auth/internal/middleware"
	"github.com/gorilla/handlers"
)

func NewRouter(c *controller.Controller) *http.ServeMux {
    router := http.NewServeMux()

    // Health Check
    router.HandleFunc("GET /ping", Fn(c.Ping))

    // Auth
    router.HandleFunc("POST /auth/register", Fn(c.Register))
    router.HandleFunc("POST /auth/login", Fn(c.Login))
    router.HandleFunc("POST /auth/logout", Fn(c.Logout))
    router.HandleFunc("GET /auth/getAuthUserID", middleware.Auth(Fn(c.GetAuthUserID), c))

    router.HandleFunc("GET /auth/google/login", Fn(c.GoogleLogin))
    router.HandleFunc("GET /auth/google/callback", Fn(c.GoogleCallback))

    return router
}

func NewCorsRouter(router *http.ServeMux, allowedOrigin string) http.Handler {
    corsHandler := handlers.CORS(
        handlers.AllowedOrigins([]string{"http://localhost:3000", "http://localhost:5173", allowedOrigin}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

    return corsHandler(router)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func Fn(fn apiFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        err := fn(w, r)
        if err != nil {
            controller.WriteJSON(w, http.StatusInternalServerError, controller.ErrMsg(err.Error()))
        }
    }
}

