package handlers

import (
	"encoding/json"
	"go-http-server/models"
	"go-http-server/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RegisterHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		hashedPassword, err := utils.HashPassword(user.Password)

		if err != nil {
			http.Error(w, "Failed To Hash Password", http.StatusInternalServerError)
			return
		}

		user.Password = hashedPassword

		if err := db.Create(&user).Error; err != nil {
			http.Error(w, "Failed To Register User", http.StatusInternalServerError)
			return
		}

		println("Register user "+user.Username+" with ID %s", user.ID)

		w.WriteHeader(http.StatusCreated)
	}
}

func LoginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqUser models.User
		if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := db.Where("username = ?", reqUser.Username).First(&user).Error; err != nil {
			http.Error(w, "User Not Found", http.StatusUnauthorized)
			return
		}

		if err := utils.CheckPassword(user.Password, reqUser.Password); err != nil {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateJWT(user.ID)
		print("token := " + token)
		if err != nil {
			http.Error(w, "Failed To Generate Token", http.StatusInternalServerError)
			return
		}

		print("Loging user "+user.Username+" with ID %s", user.ID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"token":    token,
		})
	}
}

func GetUserInfoHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")

		var user models.User
		if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
			http.Error(w, "User Not Found", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ID":       user.ID,
			"username": user.Username,
		})
	}
}

func GetUserStatsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Context().Value("userID").(uint)

		var stats models.Stats

		if err := db.First(&stats, "user_id = ?", userID).Error; err != nil {
			http.Error(w, "Stats Not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
