package db

import (
	"database/sql"
	"fmt"
	"log"
	
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// GameStat - structure for sending data to the frontend
type GameStat struct {
    ID               int    `json:"id"`
    PlayerName       string `json:"playerName"`
    Difficulty       string `json:"difficulty"`
    TimeTakenSeconds int    `json:"timeTakenSeconds"`
    SolvedAt         string `json:"solvedAt"`
}

var DB *sql.DB // Exporting a DB variable

// Init Initializes a connection to the database.
func Init(connStr string) error {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("db connection error: %w", err)
	}
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	log.Println("db connected")

	if err := goose.SetDialect("postgres"); err != nil {
        log.Fatalf("goose: failed to set dialect: %v", err)
    }

    log.Println("Running database migrations...")
    if err := goose.Up(DB, "db/migrations"); err != nil {
        log.Fatalf("goose: failed to migrate: %v", err)
    }
    log.Println("Migrations applied successfully!")

	return nil
}

// Close closes the connection to the database.
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// SaveStat saves the game result.
func SaveStat(playerName, difficulty string, timeTaken int) error {
	sqlStatement := `
		INSERT INTO game_stats (player_name, difficulty, time_taken_seconds)
		VALUES ($1, $2, $3)`
	_, err := DB.Exec(sqlStatement, playerName, difficulty, timeTaken)
	if err != nil {
		return fmt.Errorf("insert error: %w", err)
	}
	return nil
}

// GetTopStats gets the top 10 results.
func GetTopStats() ([]GameStat, error) {
	rows, err := DB.Query(`
		SELECT id, player_name, difficulty, time_taken_seconds, TO_CHAR(solved_at, 'YYYY-MM-DD HH24:MI') as solved_at
		FROM game_stats
		ORDER BY time_taken_seconds ASC
		LIMIT 10`)
	if err != nil {
		return nil, fmt.Errorf("select error: %w", err)
	}
	defer rows.Close()

	stats := make([]GameStat, 0)
	for rows.Next() {
		var stat GameStat
		if err := rows.Scan(&stat.ID, &stat.PlayerName, &stat.Difficulty, &stat.TimeTakenSeconds, &stat.SolvedAt); err != nil {
			log.Printf("scan error: %v", err)
			continue
		}
		stats = append(stats, stat)
	}

	return stats, nil
}