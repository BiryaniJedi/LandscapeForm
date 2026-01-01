package main

import (
	"encoding/json"
	"log"
	"net/http"
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
	mux := http.NewServeMux()

	mux.HandleFunc("/health", jsonHandler)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
