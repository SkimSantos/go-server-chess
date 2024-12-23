package routes

import (
	"go-http-server/handlers"
	"go-http-server/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) http.Handler {
	router := chi.NewRouter()

	// Public Routes
	router.Group(func(r chi.Router) {
		r.Post("/register", handlers.RegisterHandler(db))
		r.Post("/login", handlers.LoginHandler(db))
	})

	// Protected Routes
	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Post("/game/start", handlers.StartGameHandler(db))
		r.Post("/game/move", handlers.MakeMoveHandler(db))
		r.Post("/game/state", handlers.GetGameStateHandler(db))

		r.Get("/stats", handlers.GetUserStatsHandler(db))
	})

	return router
}
