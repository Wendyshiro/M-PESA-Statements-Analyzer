package handlers

import (
	"encoding/json"
	"log"
	"mpesa-finance/internal/auth"
	"net/http"

	"github.com/google/uuid"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type AuthHandler struct {
	authService *auth.Service
	// fot now we will use in-memory map
	users map[string]User
}

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:""`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		users:       make(map[string]User),
	}
}

func respondError(w http.ResponseWriter, message, code string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Message: message,
		Code:    code,
	})
}

func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", "INVALID_METHOD", http.StatusBadRequest)
		return
	}
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", "INVALID_JSON", http.StatusBadRequest)
		return
	}
	//Basic validation

	if req.Email == "" || req.Password == "" {
		respondError(w, "Email and password are required", "INVALID_INPUT", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8 {
		respondError(w, "Password must be at least 8 characters long", "WEAK_PASSWORD", http.StatusBadRequest)
		return
	}
	//Check is user exist
	for _, user := range h.users {
		if user.Email == req.Email {
			respondError(w, "Failed to create user", "USER_EXISTS", http.StatusConflict)
		}
	}
	//Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		respondError(w, "Failed to create user", "INTERNAL_ERROR", http.StatusInternalServerError)
	}

	//Create User
	user := User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: passwordHash,
	}
	h.users[user.ID] = user
	// Generate token
	token, err := h.authService.GenerateToken(user.ID, user.Email)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		respondError(w, "Failed to generate token", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}
	response := AuthResponse{
		Token: token,
		User:  user,
	}
	respondJSON(w, response, http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", "INVALID_METHOD", http.StatusBadRequest)
		return
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", "INVALID_JSON", http.StatusBadRequest)
	}

	//Find user
	var foundUser *User
	for _, user := range h.users {
		if user.Email == req.Email {
			foundUser = &user
			break
		}
	}
	if foundUser == nil {
		respondError(w, "Invalid email or password", "INVALID_CREDENTIALS", http.StatusUnauthorized)
		return
	}
	if err := auth.CheckPassword(foundUser.PasswordHash, req.Password); err != nil {
		respondError(w, "Invalid email or password", "INVALID_CREDENTIALS", http.StatusUnauthorized)
		return
	}
	token, err := h.authService.GenerateToken(foundUser.ID, foundUser.Email)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		respondError(w, "Failed to generate token", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}
	response := AuthResponse{
		Token: token,
		User:  *foundUser,
	}
	respondJSON(w, response, http.StatusOK)
}
