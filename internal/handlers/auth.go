package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"transaction-logger/internal/auth"
	"transaction-logger/internal/models"
)

type AuthHandler struct {
	db *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists
	_, err := models.GetUserByEmail(h.db, req.Email)
	if err == nil {
		http.Error(w, "email already in use", http.StatusBadRequest)
		return
	}

	// Create new user
	user, err := models.CreateUser(h.db, req.Email, req.Password)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "error generating token", http.StatusInternalServerError)
		return
	}

	// Return token and user info (without password)
	user.Password = ""
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user by email
	user, err := models.GetUserByEmail(h.db, req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := models.CheckPassword(user.Password, req.Password); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "error generating token", http.StatusInternalServerError)
		return
	}

	// Return token and user info (without password)
	user.Password = ""
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// AuthMiddleware verifies the JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for login and register endpoints
		if r.URL.Path == "/api/auth/register" || r.URL.Path == "/api/auth/login" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		_, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
