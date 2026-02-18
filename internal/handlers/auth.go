package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"mpesa-finance/internal/auth"
	"mpesa-finance/internal/models"
	"mpesa-finance/internal/repository"

	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *auth.Service
	userRepo    *repository.UserRepository
}

func NewAuthHandler(authService *auth.Service, userRepo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userRepo:    userRepo,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func respondError(w http.ResponseWriter, message, code string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
		"code":  code,
	})
}

func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", "INVALID_JSON", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		respondError(w, "Email & Password required", "INVALID_INPUT", http.StatusBadRequest)
		return

	}
	if len(req.Password) < 8 {
		respondError(w, "Password must be at least 8 characters", "WEAK_PASSWORD", http.StatusBadRequest)
		return
	}
	//check if user exists
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	existingUser, _ := h.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		respondError(w, "User already exists", "USER_EXISTS", http.StatusConflict)
		return
	}
	//hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		respondError(w, "Failed to create user", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}
	//create user
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: passwordHash,
	}
	if err := h.userRepo.Create(ctx, user); err != nil {
		log.Printf("Failed to create user:%v", err)
		respondError(w, "Failed to create user", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}
	//generate token
	token, err := h.authService.GenerateToken(user.ID, user.Email)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		respondError(w, "Failed to create user", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token: token,
		User:  *user,
	}
	respondJSON(w, response, http.StatusCreated)

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", "INVALID_JSON", http.StatusBadRequest)
		return
	}
	//find user
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, err := h.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("DEBUG: User not found: %v", err)
		respondError(w, "Invalid credentials", "INVALID_CREDENTIALS", http.StatusUnauthorized)
		return
	}
	//logs
	log.Printf("DEBUG: Found user: %s, hash: %s", user.Email, user.PasswordHash[:20]+"....")
	log.Printf("DEBUG: Checking password: %s", req.Password)
	//check password
	if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
		log.Printf("DEBUG: Password check FAILED! Error: %v", err)
		log.Printf("DEBUG: Hash length: %d", len(user.PasswordHash))
		respondError(w, "Invalid credentials", "INVALID_CREDENTIALS", http.StatusUnauthorized)
		return
	}
	log.Printf("DEBUG: Password check PASSED! Generating token...")
	//generate token
	token, err := h.authService.GenerateToken(user.ID, user.Email)
	if err != nil {
		log.Printf("DEBUG:Failed to generate token:%v", err)
		respondError(w, "Failed to create token", "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}
	log.Printf("DEBUG: Token generated successfully!")
	response := AuthResponse{
		Token: token,
		User:  *user,
	}
	respondJSON(w, response, http.StatusOK)
}
