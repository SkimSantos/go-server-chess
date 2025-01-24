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
		r.Handle("/*", serveStaticFiles("./static"))
		r.Post("/register", handlers.RegisterHandler(db))
		r.Post("/login", handlers.LoginHandler(db))
	})

	// Protected Routes
	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// r.Post("/game/start", handlers.StartGameHandler(db))
		// r.Post("/game/move", handlers.MakeMoveHandler(db))
		// r.Post("/game/state", handlers.GetGameStateHandler(db))

		// r.Post("/game/create", handlers.CreateGameHandler(db)) // Join now creates if no game found
		r.Post("/game/join", handlers.JoinGameHandler(db))
		r.Get("/game/get/{id}", handlers.GetGameHandler(db))
		r.Post("/game/move", handlers.SubmitMoveHandler(db))
		r.Post("/game/resign", handlers.ResignHandler(db))

		r.Get("/user/get/{id}", handlers.GetUserInfoHandler(db))
		r.Get("/user/stats", handlers.GetUserStatsHandler(db))
	})

	return router
}

func serveStaticFiles(dir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := dir + r.URL.Path
		if filePath == "./static/styles.css" {
			w.Header().Set("Content-Type", "text/css")
		}
		http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
	})
}
