package main

import (
	"go-http-server/models"
	"go-http-server/routes"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("chess-bot.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed To Connect To The Database: ", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.UsersGame{}, &models.Stats{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Cleanup games table
	if err := db.Exec("DELETE FROM users_games").Error; err != nil {
		log.Fatal("Failed to clean up users_games table:", err)
	}
	log.Println("Games table cleaned up successfully.")

	return db
}

func main() {

	db := InitDatabase()

	router := routes.SetupRoutes(db)

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
