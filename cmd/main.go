package main

import (
	"log"
	"net/http"
	"time"
	"math/rand"
	"os"

	"sudokugame/internal/db"
	"sudokugame/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("warning: .env file not found")
	}

	rand.Seed(time.Now().UnixNano())

	connStr := os.Getenv("DB_DSN")
	if connStr == "" {
		log.Fatal("DB_DSN environment variable is not set")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Fatal("APP_PORT environment variable is not set")
	}

	// DB intialization
	if err := db.Init(connStr); err != nil {
		log.Fatal("failed to init db: ", err)
	}
	defer db.Close()
	
	// HTTP-routs setting up
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/generate", handlers.GenerateHandler)
	http.HandleFunc("/stats", handlers.StatsHandler)

	log.Printf("start server http://localhost%s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}