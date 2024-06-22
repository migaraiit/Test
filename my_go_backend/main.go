package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

const (
	DB_USER     = "root"
	DB_PASSWORD = "password"
	DB_NAME     = "studyflow"
)

var db *sql.DB

func init() {
	// Initialize database connection outside main
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	var err error
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

func main() {
	defer db.Close() // Close connection on exit

	http.HandleFunc("/", getVideoURLHandler)

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getVideoURLHandler(w http.ResponseWriter, r *http.Request) {
	// Query your PostgreSQL database for the video URL
	var videoURL string
	videoID := r.URL.Query().Get("id")
	if videoID == "" {
		http.Error(w, "Missing video ID", http.StatusBadRequest)
		return
	}

	err := db.QueryRow("SELECT url FROM video_lessons WHERE id = $1", videoID).Scan(&videoURL)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Video not found", http.StatusNotFound)
		} else {
			log.Println("Error querying video URL:", err)
			http.Error(w, "Failed to fetch video URL", http.StatusInternalServerError)
		}
		return
	}

	// Return the video URL as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": videoURL})
}
