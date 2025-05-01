package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Item represents a simple data model
type Item struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Items is our in-memory database
var Items = []Item{
	{ID: "1", Name: "Item 1", CreatedAt: time.Now().Add(-24 * time.Hour)},
	{ID: "2", Name: "Item 2", CreatedAt: time.Now().Add(-12 * time.Hour)},
	{ID: "3", Name: "Item 3", CreatedAt: time.Now()},
}

// respondWithJSON writes the response as JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// errorResponse formats error responses
func errorResponse(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// healthHandler returns a simple health check response
func healthHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{
		"status": "hatt jaaa healthy buddy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// getItemsHandler returns all items
func getItemsHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, Items)
}

// getItemHandler returns a specific item by ID
func getItemHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/items/"):]

	for _, item := range Items {
		if item.ID == id {
			respondWithJSON(w, http.StatusOK, item)
			return
		}
	}

	errorResponse(w, http.StatusNotFound, "Item not found")
}

// createItemHandler adds a new item
func createItemHandler(w http.ResponseWriter, r *http.Request) {
	var item Item

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&item); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Set creation time and add a simple ID (in production, use UUIDs)
	item.CreatedAt = time.Now()
	item.ID = time.Now().Format("20060102150405")

	Items = append(Items, item)
	respondWithJSON(w, http.StatusCreated, item)
}

func main() {
	// Health check endpoint
	http.HandleFunc("/health", healthHandler)

	// API endpoints
	http.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getItemsHandler(w, r)
		case http.MethodPost:
			createItemHandler(w, r)
		default:
			errorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	// Handle single item requests
	http.HandleFunc("/api/items/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getItemHandler(w, r)
		default:
			errorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	// Start the server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
