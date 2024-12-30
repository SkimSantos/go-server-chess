package models

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Games    []Game
}

type Game struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index"`
	State     string
	Status    string
	CreatedAt time.Time
}

type UsersGame struct {
	ID          uint      `gorm:"primaryKey"`
	Player1ID   uint      `gorm:"index"` // Foreign key to User table (Player 1)
	Player2ID   uint      `gorm:"index"` // Foreign key to User table (Player 2)
	CurrentTurn uint      // User ID of the player whose turn it is
	State       string    // Board state in FEN notation
	Status      string    // "ongoing", "win", "loss", "draw"
	WinnerID    *uint     // ID of the winning player (if applicable)
	CreatedAt   time.Time // Timestamp of game creation
}

type Stats struct {
	UserID     uint `gorm:"primaryKey"`
	TotalGames int
	Wins       int
	Losses     int
	Draws      int
	MostPlayed string
}
