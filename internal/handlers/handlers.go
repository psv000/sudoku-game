package handlers

import (
	"encoding/json"
	"net/http"

	"sudokugame/internal/db"
	"sudokugame/internal/sudoku"
)

// StatsPayload - structure for receiving data from the frontend
type StatsPayload struct {
    PlayerName       string `json:"playerName"`
    Difficulty       string `json:"difficulty"`
    TimeTakenSeconds int    `json:"timeTakenSeconds"`
}

// StatsHandler routes GET and POST requests
func StatsHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getStats(w, r)
    case http.MethodPost:
        saveStats(w, r)
    default:
        http.Error(w, "method is not supported", http.StatusMethodNotAllowed)
    }
}

// saveStats saves the game result in the database
func saveStats(w http.ResponseWriter, r *http.Request) {
    var payload StatsPayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "wrong request", http.StatusBadRequest)
        return
    }

    if payload.PlayerName == "" || payload.TimeTakenSeconds <= 0 {
        http.Error(w, "player name and time must be specified", http.StatusBadRequest)
        return
    }

    if err := db.SaveStat(payload.PlayerName, payload.Difficulty, payload.TimeTakenSeconds); err != nil {
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// getStats gets top 10 results from DB
func getStats(w http.ResponseWriter, r *http.Request) {
    stats, err := db.GetTopStats()
    if err != nil {
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}

// GenerateHandler processes requests to create a new game
func GenerateHandler(w http.ResponseWriter, r *http.Request) {
    difficulty := r.URL.Query().Get("difficulty")
    if difficulty == "" {
        difficulty = "medium"
    }

    puzzle, solution := sudoku.GenerateSudoku(difficulty)
    response := sudoku.PuzzleResponse{
        Puzzle:   puzzle,
        Solution: solution,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}