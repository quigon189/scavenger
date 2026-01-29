package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"auth-service/internal/models"
	"auth-service/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	sessionName string
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		sessionName: "session_id",
	}
}

func (h *AuthHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /register", h.Register)
	router.HandleFunc("POST /login", h.Login)

	// require auth
	router.HandleFunc("GET /logout", h.RequireAuth(h.Logout))
	router.HandleFunc("POST /validate", h.RequireAuth(h.ValidateSession))
	router.HandleFunc("POST /change_password", h.RequireAuth(h.ChangePassword))
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Invalid request body: %v", err)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, "validation error", http.StatusBadRequest)
		log.Printf("Validation error: %v", err)
		return
	}

	user, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		http.Error(w, "failed to register", http.StatusBadRequest)
		log.Printf("Register error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Invalid request body: %v", err)
		return
	}

	session, sessionID, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			log.Printf("Invalid credentials: %v", err)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Printf("Internal server error: %v", err)
		return
	}

	response := &models.AuthResponse{
		User:      session.User,
		SessionID: sessionID,
		ExpiresAt: session.ExpiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := r.Context().Value("session_id").(string)
	if !ok {
		http.Error(w, "session not found", http.StatusBadRequest)
		log.Printf("session not found")
		return
	}

	if err := h.authService.Logout(r.Context(), sessionID); err != nil {
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		log.Printf("Failed to logout: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func (h *AuthHandler) ValidateSession(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := r.Context().Value("session_id").(string)
	if !ok {
		http.Error(w, "session not found", http.StatusBadRequest)
		log.Printf("session not found")
		return
	}

	session, ok := r.Context().Value("session").(*models.UserSession)
	if !ok {
		http.Error(w, "session not found", http.StatusBadRequest)
		log.Printf("session not found")
		return
	}

	response := &models.AuthResponse{
		User:      session.User,
		SessionID: sessionID,
		ExpiresAt: session.ExpiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value("session").(*models.UserSession)
	if !ok {
		http.Error(w, "session not found", http.StatusBadRequest)
		log.Printf("session not found")
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		log.Printf("invalid request body: %v", err)
		return
	}

	if err := h.authService.ChangePassword(r.Context(), session.User.ID, req.OldPassword, req.NewPassword); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Printf("bad request: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password changed successfully"})
}

func (h *AuthHandler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("X-Session-ID")
		if sessionID != "" {
			session, err := h.authService.ValidateSession(r.Context(), sessionID)
			if err == nil {
				ctx := r.Context()
				ctx = context.WithValue(ctx, "session_id", sessionID)
				ctx = context.WithValue(ctx, "session", session)

				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}
}
