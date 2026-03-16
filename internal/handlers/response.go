package handlers

import (
	"encoding/json"
	"net/http"
)

func respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)  // ← Called once here
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, message, code string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)  // ← Called once here
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
		"code":  code,
	})
}