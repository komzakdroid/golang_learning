package handlers

import (
	"dynamic-ui-backend/internal/auth"
	"dynamic-ui-backend/internal/models"
	"dynamic-ui-backend/internal/repositories"
	"dynamic-ui-backend/pkg/logger"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo *repositories.UserRepository
	logger   *logger.Logger
}

func NewAuthHandler(userRepo *repositories.UserRepository, log *logger.Logger) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, logger: log}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		h.logger.Errorw("User not found", "username", req.Username)
		h.respondError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		h.logger.Errorw("Invalid password", "username", req.Username)
		h.respondError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		h.logger.Errorw("Failed to generate token", "error", err)
		h.respondError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	expiryHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if expiryHours == 0 {
		expiryHours = 24
	}
	expiresAt := time.Now().Add(time.Hour * time.Duration(expiryHours))

	if err := h.userRepo.CreateSession(user.ID, token, expiresAt); err != nil {
		h.logger.Errorw("Failed to create session", "error", err)
	}

	response := models.LoginResponse{
		Success: true,
		Token:   token,
		User:    *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	h.logger.Infow("User logged in", "username", user.Username, "role", user.Role)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token != "" && len(token) > 7 {
		token = token[7:] // Remove "Bearer "
		h.userRepo.DeleteSession(token)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*auth.Claims)

	user, err := h.userRepo.GetByUsername(claims.Username)
	if err != nil {
		h.respondError(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    user,
	})
}

func (h *AuthHandler) respondError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Success: false,
		Error:   message,
	})
}
