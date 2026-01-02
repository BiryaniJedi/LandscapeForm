package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/BiryaniJedi/LandscapeForm-backend/internal/db"
	"github.com/joho/godotenv"
)

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Prepare the data to be sent.
	response := APIResponse{
		Status:  "success!",
		Message: "ok",
		Code:    http.StatusOK, // 200
	}

	// 2. Set the Content-Type header to "application/json".
	// This tells the client how to interpret the response body.
	w.Header().Set("Content-Type", "application/json")

	// 3. Set the HTTP status code.
	// You must set headers and status code *before* writing the body.
	w.WriteHeader(http.StatusOK)

	// 4. Encode the struct to JSON and write it to the response writer.
	// This is more memory efficient than marshalling to a byte slice first,
	// especially for large responses.
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		// Handle potential encoding errors.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	database, err := db.New()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", jsonHandler)

	if err := database.Ping(); err != nil {
		log.Printf("Error: %s\n", err)
		return
	}
	log.Println("listening on :8080, Database Connected!")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), mux))
}
