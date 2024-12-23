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

type Stats struct {
	UserID     uint `gorm:"primaryKey"`
	TotalGames int
	Wins       int
	Losses     int
	Draws      int
	MostPlayed string
}
