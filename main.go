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

	if err := db.AutoMigrate(&models.User{}, &models.Game{}, &models.Stats{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return db
}

func SeedDatabaseTemp(db *gorm.DB) {
	user := models.User{Username: "test_user", Password: "hashed_password"}
	db.FirstOrCreate(&user)

	game := models.Game{
		UserID: user.ID,
		State:  "start_position_fen",
		Status: "ongoing",
	}
	db.FirstOrCreate(&game)

	stats := models.Stats{
		UserID:     user.ID,
		TotalGames: 1,
		Wins:       0,
		Losses:     0,
		Draws:      0,
		MostPlayed: "e4",
	}
	db.FirstOrCreate(&stats)
}

func main() {

	db := InitDatabase()
	SeedDatabaseTemp(db)

	router := routes.SetupRoutes(db)

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
