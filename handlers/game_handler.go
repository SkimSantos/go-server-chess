package handlers

import (
	"encoding/json"
	"go-http-server/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/notnil/chess"
	"gorm.io/gorm"
)

func GetGameHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse game ID from URL
		gameID := chi.URLParam(r, "id")

		// Fetch the game from the database
		var game models.UsersGame
		if err := db.First(&game, "id = ?", gameID).Error; err != nil {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}

		// Return the game state as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	}
}

func CreateGameHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the ID of the current user (from AuthMiddleware context)
		player1ID := r.Context().Value("userID").(uint)

		// Initialize a new game
		game := models.UsersGame{
			Player1ID: player1ID,
			State:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", // Replace with actual starting FEN
			Status:    "pending",                                                  // Game is pending until another player joins
		}

		// Save the game to the database
		if err := db.Create(&game).Error; err != nil {
			http.Error(w, "Failed to create game", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(game)
	}
}

func JoinGameHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Fetch pending games from db and use first
		var games []models.UsersGame
		if err := db.Where("status = ?", "pending").Find(&games).Error; err != nil {
			http.Error(w, "Failed to fetch games", http.StatusInternalServerError)
			return
		}

		if len(games) == 0 {
			http.Error(w, "No games found", http.StatusInternalServerError)
			return
		}

		var req struct {
			GameID uint `json:"game_id"`
		}

		req.GameID = games[0].ID

		// Get the ID of the current user (from AuthMiddleware context)
		player2ID := r.Context().Value("userID").(uint)

		// Fetch the pending game
		var game models.UsersGame
		if err := db.First(&game, "id = ? AND status = ?", req.GameID, "pending").Error; err != nil {
			http.Error(w, "Game not found or already joined", http.StatusNotFound)
			return
		}

		// Update the game to include Player2 and start the game
		game.Player2ID = player2ID
		game.CurrentTurn = game.Player1ID // Player1 starts
		game.Status = "ongoing"

		// Save the updated game to the database
		if err := db.Save(&game).Error; err != nil {
			http.Error(w, "Failed to join game", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(game)
	}
}

func SubmitMoveHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			GameID uint   `json:"game_id"`
			Move   string `json:"move"` // e.g., "e2e4"
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		userID := r.Context().Value("userID").(uint)

		var game models.UsersGame
		if err := db.First(&game, "id = ?", req.GameID).Error; err != nil {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}

		if game.Status != "ongoing" {
			http.Error(w, "Game is not ongoing", http.StatusBadRequest)
			return
		}

		if game.CurrentTurn != userID {
			http.Error(w, "Not your turn", http.StatusForbidden)
			return
		}

		print("\n")
		print("\n")
		print(game.State)
		print("\n")
		print(req.Move)
		print("\n")
		fen, err := chess.FEN(game.State)
		if err != nil {
			http.Error(w, "Invalid game", http.StatusBadRequest)
			return
		}
		chessGame := chess.NewGame(fen, chess.UseNotation(chess.UCINotation{}))
		if err := chessGame.MoveStr(req.Move); err != nil {
			http.Error(w, "Invalid move", http.StatusBadRequest)
			return
		}

		game.State = chessGame.Position().String()
		game.CurrentTurn = togglePlayer(game.CurrentTurn, game.Player1ID, game.Player2ID)

		switch chessGame.Outcome() {
		case chess.WhiteWon:
			game.Status = "finnish"
			game.WinnerID = &game.Player1ID
			updateStats(db, game.Player1ID, game.Player1ID, false)
		case chess.BlackWon:
			game.Status = "finnish"
			game.WinnerID = &game.Player2ID
			updateStats(db, game.Player2ID, game.Player1ID, false)
		case chess.Draw:
			game.Status = "finnish"
			updateStats(db, game.Player1ID, game.Player2ID, true)
		default:
			game.Status = "ongoing"
		}

		if err := db.Save(&game).Error; err != nil {
			http.Error(w, "Failed to update game state", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	}
}

func togglePlayer(currentTurn, player1ID, player2ID uint) uint {
	if currentTurn == player1ID {
		return player2ID
	}
	return player1ID
}

func ResignHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			GameID uint `json:"game_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		userID := r.Context().Value("userID").(uint)

		var game models.UsersGame
		if err := db.First(&game, "id = ?", req.GameID).Error; err != nil {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}

		if game.Status != "ongoing" {
			if game.Status == "pending" {
				// Cleanup game
				if err := db.Exec("DELETE FROM users_games WHERE id = ?", game.ID).Error; err != nil {
					http.Error(w, "Error deleting game from database", http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(game)
				return
			}
			http.Error(w, "Game is already finished", http.StatusBadRequest)
			return
		}

		game.Status = "resigned"
		if userID == game.Player1ID {
			game.WinnerID = &game.Player2ID
			updateStats(db, game.Player2ID, game.Player1ID, false)
		} else {
			game.WinnerID = &game.Player1ID
			updateStats(db, game.Player1ID, game.Player2ID, false)
		}

		if err := db.Save(&game).Error; err != nil {
			http.Error(w, "Failed to update game state", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(game)
	}
}

func updateStats(db *gorm.DB, winnerID uint, loserID uint, isDraw bool) error {
	// Increment stats for the winner
	var winnerStats models.Stats
	db.FirstOrCreate(&winnerStats, models.Stats{UserID: winnerID})
	winnerStats.TotalGames++

	// Increment stats for the loser
	var loserStats models.Stats
	db.FirstOrCreate(&loserStats, models.Stats{UserID: loserID})
	loserStats.TotalGames++

	if isDraw {
		winnerStats.Draws++
		loserStats.Draws++
	} else {
		winnerStats.Wins++
		loserStats.Losses++
	}

	db.Save(&winnerStats)
	db.Save(&loserStats)

	return nil
}
