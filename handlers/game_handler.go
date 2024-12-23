package handlers

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

func StartGameHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("{\"message\" : \"Game Started\"}"))
	}
}

func MakeMoveHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"message\" : \"Move Made\"}"))
	}
}

func GetGameStateHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"state": "current game state"})
	}
}
