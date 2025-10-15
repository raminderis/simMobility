package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// /api/v1/test-sessions"
// /api/v1/results"
// /api/v1/reservations"

func generateReservationID() string {
	const letters = "abcdefghijklmnopqrstuvwxyz"

	code := make([]byte, 4)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}

	return "resv-" + string(code)
}

func getSessionHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Create a response object
	response := map[string]string{"message": "Sessions GET endpoint"}

	// Encode and send the JSON response
	json.NewEncoder(w).Encode(response)
}

func createSessionHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	// Create a response object
	// response := map[string]string{"message": "Sessions CREATE endpoint"}
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Add metadata
	payload["id"] = generateReservationID()
	payload["status"] = "CREATED"
	payload["created"] = time.Now().UTC().Format(time.RFC3339)
	// fmt.Println("Payload: ", payload)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payload)
}

func startSessionHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	// Create a response object
	// response := map[string]string{"message": "Sessions CREATE endpoint"}
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	payload["status"] = "STARTED"
	json.NewEncoder(w).Encode(payload)
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	// Create a response object
	response := map[string]string{"message": "Results endpoint"}
	json.NewEncoder(w).Encode(response)
}

func sessionStateHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	// Extract the session ID from the URL
	sessionID := chi.URLParam(r, "sessionID")
	// Create a response object
	response := map[string]string{
		"status":     "COMPLETED",
		"session_id": sessionID,
	}
	json.NewEncoder(w).Encode(response)
}

func deleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the session ID from the URL
	sessionID := chi.URLParam(r, "sessionID")
	w.WriteHeader(http.StatusNoContent)

	// Simulate deletion logic (e.g., remove from DB or in-memory store)
	// For now, just echo back the deleted ID
	response := map[string]string{
		"message":    "Session deleted",
		"session_id": sessionID,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func reserveSut(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	// Create a response object
	// Parse incoming JSON payload
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Add metadata
	payload["id"] = generateReservationID()
	payload["startTime"] = time.Now().UTC().Format(time.RFC3339)
	payload["created"] = time.Now().UTC().Format(time.RFC3339)
	// fmt.Println("Payload: ", payload)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payload)
}

func removeSutReservation(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	// Create a response object
	response := map[string]string{"message": "Remove SUT reservation endpoint"}
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := chi.NewRouter()
	r.Get("/api/v1/test-sessions", getSessionHandler)
	r.Post("/api/v1/test-sessions", createSessionHandler)
	r.Get("/api/v1/test-sessions/{sessionID}", sessionStateHandler)
	r.Delete("/api/v1/test-sessions/{sessionID}", deleteSessionHandler)
	r.Post("/api/v1/test-sessions/{sessionID}/execution", startSessionHandler)
	r.Get("/api/v1/results", resultsHandler)
	r.Post("/api/v1/reservations", reserveSut)
	r.Delete("/api/v1/reservations/{reservationID}", removeSutReservation)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
